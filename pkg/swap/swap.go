// Package swap chooses the transport based on a given payload threshold
package swap

import (
	"errors"

	"github.com/bign8/msg"
)

var _ msg.Transport = (*Transport)(nil)

// New constructs a new Transport
func New(pubsub, blob msg.Transport, change int) *Transport {
	if change <= 0 {
		panic("WTF are you doing?")
	}
	return &Transport{
		pubsub: pubsub,
		blob:   blob,
		change: change,
	}
}

// Transport is a msg.Transport
type Transport struct {
	pubsub, blob msg.Transport
	change       int
}

// Kill closes the transport
func (t *Transport) Kill() error {
	if err := t.blob.Kill(); err != nil {
		return err
	}
	return t.pubsub.Kill()
}

// Wait returns when transport closes
func (t *Transport) Wait() <-chan error {
	c := make(chan error, 1)
	c <- errors.New("TODO")
	return c
}

// Able verifies the transport is open
func (t *Transport) Able() bool {
	return t.blob.Able() && t.pubsub.Able()
}

// Open starts up a transport
func (t *Transport) Open() error {
	return errors.New("TODO")
}

// Push does a one way transaction
func (t *Transport) Push(ctx msg.Context, subject string, data []byte) error {
	if len(data) > t.change {
		return t.blob.Push(ctx, subject, data)
	}
	return t.pubsub.Push(ctx, subject, data)
}

// Recv receives data (blocking call if suppported)
func (t *Transport) Recv(fn func(string, []byte)) error {
	return nil
}

// Send does an RPC stype round trip
func (t *Transport) Send(ctx msg.Context, subject string, data []byte) ([]byte, error) {
	if len(data) > t.change {
		return t.blob.Send(ctx, subject, data)
	}
	return t.pubsub.Send(ctx, subject, data)
}
