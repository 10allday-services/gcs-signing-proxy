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
	obj := proxy.bucket.Object(path)

	// If the reader can't find anything, it's probably a 404
	reader, err := obj.NewReader(ctx)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	defer reader.Close()

	// FIXME(willkg): do we need to set other headers?
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Proxy all the data from the reader to the response
	_, err = io.Copy(w, reader)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		panic(err)
	}
}
