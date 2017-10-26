// Package msg contains core messaging constructs.
//
// These constructs are specifically designed to be transpiled with gopherjs.
package msg

import (
	"errors"
	"time"
)

const (
	rpcReqPrefix = "rpc.req."
	rpcResPrefix = "rpc.res."
)

// Common package errors
var (
	ErrClosed = errors.New("CLOSED")
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
	Open() error // blocking call
	Close() error
}

// TRPC supports a full RPC structure
type TRPC interface {
	Transport
	Recv(func(string, []byte)) error // blocking call
	Send(string, []byte) error
}

// TPubSub supports a publish/subscribe model of comunications
type TPubSub interface {
	Transport
	Publish(name string, data []byte) error
	Subscribe(name string, cb func([]byte)) (Stream, error)
}

// Stream lets you listen to multiple messages on a socket
type Stream func() error

// Close from a process stream
func (s Stream) Close() error { return s() }
