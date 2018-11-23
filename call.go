package msg

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
)

type contextKey int

const (
	contextKeyCaller = contextKey(iota) // Caller
)

// With attaches various things to a given context
func With(ctx Context, subject interface{}) context.Context {
	var key contextKey
	switch subject.(type) {
	case Caller:
		key = contextKeyCaller
	default:
		panic("with: unknown subject interface")
	}
	return context.WithValue(ctx, key, subject)
}

// Caller is the base call type
type Caller func(ctx context.Context, service, method string, in io.Reader, out io.Writer) error

// Call creates a service call to an external service on your behalf.
func Call(ctx context.Context, service, method string, in io.Reader, out io.Writer) error {
	caller, ok := ctx.Value(contextKeyCaller).(Caller)
	if !ok {
		return errors.New("call: no caller on context")
	}
	return caller(ctx, service, method, in, out)
}

// Main is the parent function to be invoked by your main process.
//
// Allows the injection of the primary caller into http handlers
func Main(ctx context.Context) {
	port := ":" + os.Getenv("PORT")
	http.ListenAndServe(port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.DefaultServeMux.ServeHTTP(w, r.WithContext(ctx))
	}))
}
