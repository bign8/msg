// Package nats holds a nats implementation of binding
package nats

import "github.com/bign8/msg"

var _ msg.Transport = (*Transport)(nil)

// Transport is a nats based transport
type Transport struct {
}

// Close kills the transport
func (t *Transport) Close() error {
	return nil
}

// IsOpen verifies the transport is open
func (t *Transport) IsOpen() bool {
	return false
}

// Open starts up a transport
func (t *Transport) Open() error {
	return nil
}

// Push does a one way transaction
func (t *Transport) Push(subject string, data []byte) error {
	return nil
}

// Recv receives data (blocking call if suppported)
func (t *Transport) Recv(fn func(string, []byte)) error {
	return nil
}

// Send does an RPC stype round trip
func (t *Transport) Send(subject string, data []byte) ([]byte, error) {
	return nil, nil
}
