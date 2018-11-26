// Package http provides an http tranport
package http

import (
	"context"
	"net/http"
	"os"
)

// Main is the parent function to be invoked by your main process.
//
// Allows the injection of the primary caller into http handlers
func Main(ctx context.Context) {
	port := ":" + os.Getenv("PORT")
	http.ListenAndServe(port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.DefaultServeMux.ServeHTTP(w, r.WithContext(ctx))
	}))
}
