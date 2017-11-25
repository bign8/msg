package msg

import (
	"bytes"
	"context"
	"errors"
	"testing"
)

func TestGenID(t *testing.T) {
	var c *Conn
	out := c.genID()
	if l := len(out); l != 8 {
		t.Errorf("Expected 8; Got %d; %q", l, out)
	}
}

type msg struct {
	title string
	body  []byte
}

func newTranz() Transport {
	return &tranz{
		tunnel: make(chan *msg, 1),
	}
}

type tranz struct {
	tunnel chan *msg
}

func (t *tranz) Open() error { return nil }
func (t *tranz) Able() bool  { return true }
func (t *tranz) Kill() error { return nil }
func (t *tranz) Wait() <-chan error {
	c := make(chan error, 1)
	c <- nil
	return c
}
func (t *tranz) Push(ctx Context, subject string, data []byte) error {
	t.tunnel <- &msg{title: subject, body: data}
	return nil
}
func (t *tranz) Recv(fn func(string, []byte)) error {
	for m := range t.tunnel {
		fn(m.title, m.body)
	}
	return nil
}
func (t *tranz) Send(ctx Context, subject string, data []byte) ([]byte, error) {
	return nil, errors.New("nope")
}

func TestConn(t *testing.T) {
	c := New(newTranz())
	var called bool
	c.Handle("hello", func(ctx Context, bits []byte) ([]byte, error) {
		called = true
		return bits, nil
	})
	bits, err := c.Request(context.TODO(), "hello", []byte("nate"))
	if err != nil {
		t.Fatalf("Didn't expect error: %q", err)
	}
	if bytes.Compare(bits, []byte("nate")) != 0 {
		t.Errorf("Wanted %q; Received %q", []byte("nate"), bits)
	}
}
