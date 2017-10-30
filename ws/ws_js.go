package ws

import (
	"github.com/bign8/msg"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket/websocketjs"
)

var _ msg.Transport = (*Transport)(nil)

// New constructs a new Transport
func New(addr string) *Transport {
	return &Transport{
		addr: addr,
	}
}

// Transport is a msg.Transport
type Transport struct {
	ws   *websocketjs.WebSocket
	addr string     // address of connection
	attp int        // attempt #
	ch   chan error // closes when transport does
}

// Kill closes the transport
func (t *Transport) Kill() error {
	if !t.Able() {
		print("Transport is already closed")
		return nil
	}
	err := t.ws.Close()
	if err == nil {
		err = <-t.ch
	}
	return err
}

// Able verifies the transport is open
func (t *Transport) Able() bool {
	return t.ws != nil && t.ch != nil
}

// Open starts up a transport
func (t *Transport) Open() (err error) {
	if t.Able() {
		print("Transport is already open")
		return nil
	}
	t.ws, err = websocketjs.New(t.addr)
	if err == nil {
		t.ch = make(chan error, 1)
		t.ws.AddEventListener("msg", false, t.onMsg)
		t.ws.AddEventListener("err", false, t.onErr)
		t.ws.AddEventListener("open", false, t.onOpen)
		t.ws.AddEventListener("close", false, t.onClose)
	}
	return err
}

// Wait returns when transport closes
func (t *Transport) Wait() <-chan error {
	return t.ch
}

func (t *Transport) onOpen(obj *js.Object) {
	t.attp = 0
}

func (t *Transport) onClose(obj *js.Object) {
	t.ws.Close()
	ch := t.ch
	t.ch = nil
	close(ch)
}

func (t *Transport) onMsg(obj *js.Object) {

}

func (t *Transport) onErr(obj *js.Object) {
	t.onErr(obj)
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
