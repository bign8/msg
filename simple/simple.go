package simple

import (
	"context"
	"errors"
	"io"
)

// Call is the fundamental API between product and infrastructure.
// 
// Usage:
//   RPC: in != nil && out != nil
//   PUB: in != nil && out == nil
//   SUB: in == nil && out != nil && (_, ok := out.(Subscription); ok)
//
// For RPCs, Call serializes all the bytes from `in` and transmits them to the server.
//   Based on the response message, the Call can return an error on failure, or hydrate the `out` object for success.
//
// For PUB the `in` is handled like an RPC, and an error is returned if the message failed to reach a messaging service.
//
// For SUB a subscription is constructed for a particular channel and is held open until the Call context is Cancled.
func Call(ctx context.Context, service, method string, in io.Reader, out io.Writer) error {
	switch {
		case in != nil && out != nil:
			return rpc(ctx, service, method, in, out)
		case in != nil && out == nil:
			return pub(ctx, service, method, in)
		case s, ok := out.(Subscription); in == nil && ok:
			return sub(ctx, service, method, s)
	}
	return errors.New("simple.Call: un-supported argument configuration")
}

// Subscription is an interface that must be implemented when subscribing to streams.
type Subscription struct {
	io.Writer

	// Reset is called by Call once the full object has been written.
  // This should clear the current object frame and prepare for the next incomming message.
	Reset() error
}

func rpc(ctx context.Context, service, method string, in io.Reader, out io.Writer) error {
	return errors.New("simple.RPC: TODO")
}

func pub(ctx context.Context, service, method string, in io.Reader) error {
	return errors.New("simple.PUB: TODO")
}

func sub(ctx context.Context, service, method string, s Subscription) error {
	return errors.New("simple.SUB: TODO")
}

