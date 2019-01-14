[![Go Report Card](https://goreportcard.com/badge/github.com/mozilla-services/gcp-signing-proxy)](https://goreportcard.com/report/github.com/mozilla-services/gcp-signing-proxy)

# gcp-signing-proxy

Proxies incoming HTTP GET requests to signed GCS object retrieval.

## Usage

Images are available from Docker Hub:

```
$ docker pull mozilla/gcp-signing-proxy
```

You'll need a Google Compute Platform service account key JSON file for a service account that has view access to the Google Cloud Storage bucket you're using. You can create a service account key in the GCP console.

Then run gcp-signing-proxy like this:

```
$ docker run -v CREDENTIALSFILE:/service_account_key.json -p 8000:8000 mozilla/gcp-signing-proxy:latest
```

replacing `CREDENTIALSFILE` with the filename of your credentials file.

This mounts your service account key JSON file as `/service_account_key.json` in the container and exposes the 8000 port to the host.

The signing proxy listens on `0.0.0.0:8000` by default, which means that it will be exposed to the world _if you expose that port externally_.

## Configuration

The signing proxy is configured via environment variables with the prefix `SIGNING_PROXY_`. The [config struct](https://github.com/mozilla-services/gcp-signing-proxy/blob/master/main.go#L83-L92) has details on default values and variable types. Implementation by Kelsey Hightower's [envconfig](github.com/kelseyhightower/envconfig).

Available environment variables:

    - SIGNING_PROXY_LOG_REQUESTS
        type: bool
        description: enable logging of request method and path
        default: "true"
    - SIGNING_PROXY_STATSD
        type: bool
        description: enable statsd reporting
        default: "true"
    - SIGNING_PROXY_STATSD_LISTEN
        type: string
        description: address to send statsd metrics to
        default: "127.0.0.1:8125"
    - SIGNING_PROXY_STATSD_NAMESPACE
        type: string
        description: prefix for statsd metrics. "." is appended as a separator.
        default: "SIGNING_PROXY"
    - SIGNING_PROXY_LISTEN
        type: string
        description: address for the proxy to listen on
        default: "0.0.0.0:8000"
    - SIGNING_PROXY_BUCKET
        type: string
        description: the GCS bucket
        default: ""


## Development/hacking

Requirements:

* docker
* make
* [`dep`](https://golang.github.io/dep/)
* [`gox`](https://github.com/mitchellh/gox)

To build binary and Docker image, do:

```
$ make build
```

To sync `Gopkg.lock` and vendored packages:

```
$ dep ensure
```

You'll need a Google Cloud Platform service account key JSON file for a service account that has view access to the Google Cloud Storage bucket you're using. You can create a service account key in the GCP console.

You'll need to create a `my.env` env file with the bucket name in it like this:

```
SIGNING_PROXY_BUCKET=mybucket
```

You can set other configuration in that file, too, or use the defaults.

Then run gcp-signing-proxy like this:

```
$ make run
```

You can also set the `CREDS` environment variable to the name of your service account key JSON file if it's not named `service_account_key.json`.

`CTRL-c` will stop the service.

After making big changes, update the `version` const in `main.go`.
