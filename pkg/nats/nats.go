// Package nats holds a nats implementation of binding
package nats

import msg "github.com/bign8/msg/pkg"

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
func (t *Transport) Push(ctx msg.ContextOld, data *msg.Msg) error {
	return nil
}

// Recv receives data (blocking call if suppported)
func (t *Transport) Recv(fn func(*msg.Msg)) error {
	return nil
}

// Send does an RPC stype round trip
func (t *Transport) Send(ctx msg.ContextOld, data *msg.Msg) (*msg.Msg, error) {
	return nil, nil
}
