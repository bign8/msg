package msg

// Context is a copy of context.Context but only the methods that are needed.
type Context interface {
	// Deadline removed because time.Time doesn't transpile well.
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

type contextKey int

const (
	contextKeyCaller = contextKey(iota) // Caller
)

// WithCaller attaches various things to a given context
func WithCaller(ctx Context, val Callable) Context {
	return WithValue(ctx, contextKeyCaller, val)
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

// Simple contexts exposed by package
var (
	Background Context = emptyCtx(0)
	TODO       Context = emptyCtx(1)
)

type emptyCtx int

func (emptyCtx) Done() <-chan struct{}         { return nil }
func (emptyCtx) Err() error                    { return nil }
func (emptyCtx) Value(interface{}) interface{} { return nil }
