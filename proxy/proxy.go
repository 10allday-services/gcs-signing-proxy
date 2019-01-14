// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package proxy

import (
	"cloud.google.com/go/storage"

	"context"
	"io"
	"net/http"
)

// New creates a Handler using the input destination and client
func New(bucket *storage.BucketHandle) (*Handler, error) {
	return &Handler{bucket}, nil
}

type Handler struct {
	bucket *storage.BucketHandle
}

// Satisfies http.Handler
// per https://golang.org/pkg/net/http/#Handler
// the server will recover panic() and log a stack trace
func (proxy Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	ctx := context.Background()

	// FIXME(willkg): might need to use RawPath here or EscapedPath
	url := *req.URL
	path := url.Path
	// Need to drop the / at the beginning of the path
	obj := proxy.bucket.Object(path[1:])

	// Get the content type; if it errors, it's probably a 404
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		http.NotFound(w, req)
		return
	}

	contentType := attrs.ContentType

	reader, err := obj.NewReader(ctx)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		panic(err)
	}
	defer reader.Close()

	// FIXME(willkg): do we need to set other headers?
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)

	// Proxy all the data from the reader to the response
	_, err = io.Copy(w, reader)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		panic(err)
	}
}
