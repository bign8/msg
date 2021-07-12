package msg

import "errors"

// Reader is a copy of io.Reader.
type Reader interface {
	Read(p []byte) (n int, err error)
}

// Writer is a copy of io.Writer.
type Writer interface {
	Write(p []byte) (n int, err error)
}

// Callable is the base call type
type Callable func(ctx Context, service, method string, in Reader, out Writer) error

// Call creates a service call to an external service on your behalf.
func Call(ctx Context, service, method string, in Reader, out Writer) error {
	caller, ok := ctx.Value(contextKeyCaller).(Callable)
	if !ok {
		return errors.New("call: no caller on context")
	}
	return caller(ctx, service, method, in, out)
}

type ServiceMethod struct {
	Service string
	Method  string
}

// Register yourself as a handler of a thing.
func Router(callers map[ServiceMethod]Callable) Callable {
	return callTree(callers).call
}

type callTree map[ServiceMethod]Callable

func (tree callTree) call(ctx Context, service, method string, in Reader, out Writer) error {
	for key, callable := range tree {
		if key.Service == service && key.Method == method {
			return callable(ctx, service, method, in, out)
		}
	}
	return errors.New(`Unprocessable entity`)
}
