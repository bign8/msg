package client

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket/websocketjs"

	"github.com/bign8/msg"
)

const (
	failover = 1024 * 1024 * 4 // TODO: benchmark
)

type network struct {
	loc string
	ws  *websocketjs.WebSocket
	// TODO: rate limiters
}

// New constructs a new msg.Caller.
func New(server string) msg.Caller {
	return (&network{}).call
}

func (net *network) call(ctx msg.Context, service, method string, in msg.Reader, out msg.Writer) error {

	// read data into memory bufffer
	bits, err := readAll(in, 512)
	if err != nil {
		return err
	}

	// Choose which transport base on message size
	if len(bits) > failover {
		bits, err = net.sendHTTP(ctx, service, method, bits)
	} else {
		bits, err = net.sendWS(ctx, service, method, bits)
	}
	if err != nil {
		return err
	}

	// Write the data back into a struct
	_, err = out.Write(bits)
	return err
}

func (net *network) newWS() error {
	ws, err := websocketjs.New("/api/ws")
	if err != nil {
		return err
	}
	open := make(chan error, 1)
	ws.AddEventListener("open", true, func(o *js.Object) {
		println("TODO(ws): open", o)
		open <- nil
	})
	ws.AddEventListener("error", true, func(o *js.Object) {
		println("TODO(ws): error", o)
	})
	ws.AddEventListener("close", true, func(o *js.Object) {
		println("TODO(ws): close", o)
	})
	ws.AddEventListener("message", true, func(o *js.Object) {
		println("TODO(ws): message", o)
	})
	net.ws = ws
	return <-open
}

func (net *network) sendWS(ctx msg.Context, service, method string, in []byte) ([]byte, error) {

	// Lazy initiate websocket (TODO: attempt in constructor)
	if net.ws == nil {
		if err := net.newWS(); err != nil {
			println("Network(newWS): failure", err)
			return net.sendHTTP(ctx, service, method, in)
		}
	}

	panic("TODO")
}

func err2chan(out chan<- error, cb func() error) func() {
	return func() { out <- cb() }
}

func (net *network) sendHTTP(ctx msg.Context, service, method string, bits []byte) (res []byte, err error) {
	done := make(chan error, 1)
	req := js.Global.Get("XMLHttpRequest").New()
	req.Call("open", "POST", "/api/"+service+"/"+method, true /* async */)
	req.Call("setRequestHeader", "Content-Type", "application/octet-stream")
	req.Set("onload", err2chan(done, func() error {
		if req.Get("status").Int() != 200 {
			return errors.New("drive.Network(" + req.Get("statusText").String() + "): " + req.Get("response").String())
		}
		data := js.Global.Get("Uint8Array").New(req.Get("response"))
		res = data.Interface().([]byte)
		return nil
	}))
	req.Set("onerror", func() {
		print("error", req)
		done <- errors.New(req.Get("response").String())
	})
	req.Set("responseType", "arraybuffer")
	req.Call("send", bits)
	return res, <-done
}

// copy of ioutil.ReadAll and buffer.ReadFrom
func readAll(r msg.Reader, capacity int) (b []byte, err error) {
	for {
		// grow
		i := len(b)
		b = append(b, make([]byte, capacity)...)
		b = b[:i]

		// regulary buf.ReadFrom loop
		m, e := r.Read(b[i:cap(b)])
		if m < 0 {
			panic("negative read")
		}
		b = b[:i+m]
		if e != nil && e.Error() == "EOF" {
			return b, nil // e is EOF, so return nil explicitly
		}
		if e != nil {
			return b, e
		}
	}
}
