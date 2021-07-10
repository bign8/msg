package simple

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	apiPath = "/rpc_http"
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
	if s, ok := out.(Subscription); in == nil && ok {
		return sub(ctx, service, method, s)
	}
	if in != nil && out == nil {
		return pub(ctx, service, method, in)
	}
	if in != nil && out != nil {
		return rpc(ctx, service, method, in, out)
	}
	return errors.New("simple.Call: un-supported argument configuration")
}

// Subscription is an interface that must be implemented when subscribing to streams.
type Subscription interface {
	io.Writer

	// Reset is called by Call once the full object has been written.
	// This should clear the current object frame and prepare for the next incomming message.
	Reset() error
}

func rpc(ctx context.Context, service, method string, in io.Reader, out io.Writer) error {
	req, err := http.NewRequest(http.MethodPost, os.Getenv("API_HOST")+apiPath, in)
	if err != nil {
		return err
	}
	// TODO: send this data encoded in the message
	req.Header.Add("X-BIGN8-SERVICE", service)
	req.Header.Add("X-BIGN8-METHOD", method)
	res, err := http.DefaultClient.Do(req.WithContext(ctx)) // TODO: use client that limits outgoing dials to 100 per process
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		body, _ := ioutil.ReadAll(res.Body)
		return errors.New("simple.RPC: unknown failure: " + string(body))
	}
	// TODO: decode possible errors from response body
	_, err = io.Copy(out, res.Body)
	return err
}

func pub(ctx context.Context, service, method string, in io.Reader) error {
	return errors.New("simple.PUB: TODO")
}

func sub(ctx context.Context, service, method string, s Subscription) error {
	return errors.New("simple.SUB: TODO")
}
