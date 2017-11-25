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

func newTranz() Transport {
	return &tranz{
		tunnel: make(chan *Msg, 1),
	}
}

type tranz struct {
	tunnel chan *Msg
}

func (t *tranz) Open() error                              { return nil }
func (t *tranz) Able() bool                               { return true }
func (t *tranz) Kill() error                              { return nil }
func (t *tranz) Send(ctx Context, msg *Msg) (*Msg, error) { return nil, errors.New("nope") }
func (t *tranz) Wait() <-chan error {
	c := make(chan error, 1)
	c <- nil
	return c
}
func (t *tranz) Push(ctx Context, msg *Msg) error {
	t.tunnel <- msg
	return nil
}
func (t *tranz) Recv(fn func(*Msg)) error {
	for m := range t.tunnel {
		m.Title = m.Reply // make it a reply
		m.Reply = ""
		fn(m)
	}
	return nil
}

func TestConn(t *testing.T) {
	c := New(newTranz())
	err := c.Open()
	if err != nil {
		t.Fatalf("Didn't expect open error: %q", err)
	}
	res, err := c.Request(context.TODO(), &Msg{Title: "hello", Body: []byte("nate")})
	if err != nil {
		t.Fatalf("Didn't expect request error: %q", err)
	}
	if bytes.Compare(res.Body, []byte("nate")) != 0 {
		t.Errorf("Wanted %q; Received %q", []byte("nate"), res.Body)
	}
}
