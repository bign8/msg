// Package claimcheck chooses the transport based on a given payload threshold.
// https://akfpartners.com/growth-blog/claim-check-pattern
package claimcheck

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
func (t *Transport) Push(ctx msg.ContextOld, data *msg.Msg) error {
	if len(data.Body) > t.change {
		return t.blob.Push(ctx, data)
	}
	return t.pubsub.Push(ctx, data)
}

// Recv receives data (blocking call if suppported)
func (t *Transport) Recv(fn func(*msg.Msg)) error {
	return nil
}

// Send does an RPC stype round trip
func (t *Transport) Send(ctx msg.ContextOld, data *msg.Msg) (*msg.Msg, error) {
	if len(data.Body) > t.change {
		return t.blob.Send(ctx, data)
	}
	return t.pubsub.Send(ctx, data)
}
