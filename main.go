package main

import (
	"cloud.google.com/go/storage"
	"google.golang.org/api/option"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/kelseyhightower/envconfig"

	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/mozilla-services/gcp-signing-proxy/proxy"
)

const (
	appNamespace = "SIGNING_PROXY"
	version      = "1.0.0"
)

var (
	statsdClient *statsd.Client
	httpClient   *http.Client
	pool         *x509.CertPool
)

// get CA certs for our http.Client
func init() {
	// cacert.pem is a runtime dependency!
	bs, err := ioutil.ReadFile("/cacert.pem")
	if err != nil {
		log.Fatal(err.Error())
	}

	pool = x509.NewCertPool()
	pool.AppendCertsFromPEM(bs)

	// default http client with a timeout
	httpClient = &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: pool},
		},
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method + " " + r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func statsdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := statsdClient.Incr("requests", []string{}, 1.0)
		if err != nil {
			log.Println(err.Error())
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("Starting gcp-signing-proxy....")
	config := struct {
		LogRequests     bool   `default:"true" split_words:"true"`
		Statsd          bool   `default:"true"`
		StatsdListen    string `default:"127.0.0.1:8125" split_words:"true"`
		StatsdNamespace string `default:""`
		Listen          string `default:"0.0.0.0:8000"`
		Service         string `default:"storage"`
		Region          string `default:"us-east1"`
		Destination     string `default:"https://www.googleapis.com/storage/v1/b"`
	}{}

	// load envconfig
	err := envconfig.Process(appNamespace, &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Loaded config")

	// *url.URL for convenience
	destinationURL, err := url.Parse(config.Destination)
	if err != nil {
		log.Fatal("Could not parse destination URL: " + err.Error())
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("/service_account_key.json"))
	if err != nil {
		log.Fatal("Could not get credentials: " + err.Error())
	}
	_ = client // Use the client.
	log.Println("Built storage client")

	// FIXME(willkg): get the credentials, create signing client, create
	// proxy using signing client

	// create proxy using signing http client
	// FIXME add signingClient as second arg
	prxy, err := proxy.New(destinationURL, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	var handler http.Handler
	handler = prxy

	// wrap handler for logging
	if config.LogRequests {
		handler = loggingMiddleware(handler)
	}

	// wrap handler for statsd
	if config.Statsd {
		statsdClient, err := statsd.New(config.StatsdListen)
		if err != nil {
			log.Fatal(err.Error())
		}
		// prepended to metrics
		if config.StatsdNamespace == "" {
			statsdClient.Namespace = appNamespace + "."
		} else {
			statsdClient.Namespace = config.StatsdNamespace + "."
		}
		// statsdClient.Tags = append(statsdClient.Tags, ec2tags...)
		handler = statsdMiddleware(handler)
	}

	log.Println("Starting " + appNamespace)

	// sane default timeouts
	srv := &http.Server{
		Addr:         config.Listen,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      handler,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
