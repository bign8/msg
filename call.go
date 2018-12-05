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
