// Package nats holds a nats implementation of binding
package nats

import "github.com/bign8/msg"

var _ msg.Transport = (*Transport)(nil)

// New constructs a new Transport
func New() *Transport {
	return &Transport{}
}

// Transport is a msg.Transport
type Transport struct {
}

// Kill closes the transport
func (t *Transport) Kill() error {
	return nil
}

// Wait returns when transport closes
func (t *Transport) Wait() <-chan error {
	return nil
}

// Able verifies the transport is open
func (t *Transport) Able() bool {
	return false
}

// Open starts up a transport
func (t *Transport) Open() error {
	return nil
}

// Push does a one way transaction
func (t *Transport) Push(ctx msg.Context, subject string, data []byte) error {
	return nil
}

// Recv receives data (blocking call if suppported)
func (t *Transport) Recv(fn func(string, []byte)) error {
	return nil
}

// Send does an RPC stype round trip
func (t *Transport) Send(ctx msg.Context, subject string, data []byte) ([]byte, error) {
	return nil, nil
}
