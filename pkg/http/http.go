// Package http provides an http tranport
package http

import (
	"net/http"
	"os"
	"strings"

	"github.com/bign8/msg"
)

// Main is the parent function to be invoked by your main process.
//
// Allows the injection of the primary caller into http handlers
func Main(callable msg.Callable) {
	port := ":" + os.Getenv("PORT")
	http.ListenAndServe(port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: check if I'm able to handle the request!
		if strings.HasPrefix(r.URL.Path, "/_n8/rpc/") {
			// strings.SplitAfterN(r.URL.Path, "/", 2)
			// TODO: parse service method out of url
			err := callable(r.Context(), `todo`, `todo`, nil, nil)
			if err != nil {
				panic(err)
				// TODO: write error
			}
			// TODO: write ok
			http.Error(w, `todo`, http.StatusNotImplemented)
			return
		}
		// ctx := msg.WithCaller(r.Context(), callable)
		http.DefaultServeMux.ServeHTTP(w, r) //.WithContext(ctx))
	}))
}
