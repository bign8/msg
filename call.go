package msg

import "errors"

type contextKey int

const (
	contextKeyCaller = contextKey(iota) // Caller
)

// WithCaller attaches various things to a given context
func WithCaller(ctx Context, val Caller) Context {
	return WithValue(ctx, contextKeyCaller, val)
}

// Reader is a copy of io.Reader.
type Reader interface {
	Read(p []byte) (n int, err error)
}

// Writer is a copy of io.Writer.
type Writer interface {
	Write(p []byte) (n int, err error)
}

// Context is a copy of context.Context but only the methods that are needed.
type Context interface {
	// Deadline removed because time.Time doesn't transpile well.
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

// WithValue is a copy of context.Context (with reflect Comparable check removed).
func WithValue(parent Context, key, val interface{}) Context {
	if key == nil {
		panic("nil key")
	}
	return &valueCtx{parent, key, val}
}

// A valueCtx carries a key-value pair. It implements Value for that key and
// delegates all other calls to the embedded Context.
type valueCtx struct {
	Context
	key, val interface{}
}

func (c *valueCtx) Value(key interface{}) interface{} {
	if c.key == key {
		return c.val
	}
	return c.Context.Value(key)
}

// Caller is the base call type
type Caller func(ctx Context, service, method string, in Reader, out Writer) error

// Call creates a service call to an external service on your behalf.
func Call(ctx Context, service, method string, in Reader, out Writer) error {
	caller, ok := ctx.Value(contextKeyCaller).(Caller)
	if !ok {
		return errors.New("call: no caller on context")
	}
	return caller(ctx, service, method, in, out)
}

// // Main is the parent function to be invoked by your main process.
// //
// // Allows the injection of the primary caller into http handlers
// func Main(ctx Context) {
// 	port := ":" + os.Getenv("PORT")
// 	http.ListenAndServe(port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		http.DefaultServeMux.ServeHTTP(w, r.WithContext(ctx))
// 	}))
// }
