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

// Msg is the core data struct for transporting messages
type Msg struct {
	Title string
	Reply string
	Body  []byte
}

// Handler is a message handler
type Handler func(*Msg) error

// Context mirrors context.Context (but with fewer imports)
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

// Transport is the core communication interface we will communicate over
type Transport interface {
	Open() error                      // start the given transport
	Able() bool                       // is this transport open
	Kill() error                      // close down this transport
	Wait() <-chan error               // channel is closed when transport is closed
	Recv(func(*Msg)) error            // blocking call - when data is received
	Send(Context, *Msg) (*Msg, error) // Send some data
	Push(Context, *Msg) error         // send a message one direction
}

// Stream lets you listen to multiple messages on a socket
type Stream func() error

// Close from a process stream
func (s Stream) Close() error { return s() }
