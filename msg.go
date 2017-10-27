// Package msg contains core messaging constructs.
//
// These constructs are specifically designed to be transpiled with gopherjs.
package msg

import (
	"errors"
	"time"
)

const (
	rpcReqPrefix = "rpc.req." // rpc.req.<service>.<version>.<function>
	rpcResPrefix = "rpc.res."
)

// Common package errors
var (
	ErrClosed = errors.New("CLOSED")

	errNotSupported = errors.New("not supported")
)

// Context mirrors context.Context (but with fewer imports)
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

// Transport is the core communication interface we will communicate over
type Transport interface {
	Open() error                                  // start the given transport
	Able() bool                                   // is this transport open
	Kill() error                                  // close down this transport
	Wait() <-chan error                           // channel is closed when transport is closed
	Recv(func(string, []byte)) error              // blocking call
	Send(Context, string, []byte) ([]byte, error) // Send some data
	Push(Context, string, []byte) error           // send a message one direction
}

// Stream lets you listen to multiple messages on a socket
type Stream func() error

// Close from a process stream
func (s Stream) Close() error { return s() }
