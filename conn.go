package msg

import (
	"errors"
	"time"
)

// New builds a new Sock that performs automatic retires on connection failures
func New(trans Transport, opts ...Option) *Conn {
	conn := &Conn{
		hands: make(map[string]func([]byte), 1),
		trans: trans,
		close: make(chan chan error, 1),
	}
	for _, o := range opts {
		o(conn)
	}
	return conn
}

// Option is a particular option you can apply to a connection
type Option func(*Conn)

// Conn is a core managed connection for all communication
type Conn struct {
	hands map[string]func([]byte) // array of active handlers
	trans Transport               // currently active transport

	// TODO: think about removing these guys
	err   error
	att   int             // retry attempt
	close chan chan error // closing
}

// Open starts a connection
func (s *Conn) Open() error {
	err := s.trans.Open()
	if err != nil {
		return err
	}
	go s.open()
	return nil
}

// open manages the connection
func (s *Conn) open() {
	for {
		s.err = s.trans.Recv(s.recv)
		if s.err == errNotSupported {
			s.err = nil
			return
		}

		// Something failed, delay and try connecting again
		var delay int
		delay, s.att = RetryDelay(s.att)
		time.Sleep(time.Duration(delay))
		s.trans.Open() // TODO: handle errors here
	}
}

func (s *Conn) recv(msg *Msg) {
	s.att = 0
	s.err = nil
	fn, ok := s.hands[msg.Title]
	if !ok {
		print("Unsupported function:" + msg.Title)
		return
	}
	fn(msg.Body)
}

func (s *Conn) genID() string {
	return "01234567" // NOTE: should always be len(8) so decoding is consistent
}

// Request executes an RPC
func (s *Conn) Request(ctx Context, name string, data []byte) ([]byte, error) {
	if s.err != nil || s.trans == nil {
		return nil, s.err
	}
	if !s.trans.Able() {
		return nil, ErrClosed
	}
	resc := make(chan []byte, 1)
	defer close(resc) // yay memory leaks

	// Generate response subscription
	reply := s.genID()
	sub, err := s.Subscribe(rpcResPrefix+reply, func(res []byte) { resc <- res })
	if err != nil {
		return nil, err
	}
	defer sub.Close() // TODO: log error here

	// Publish request to the cloud
	if err := s.trans.Push(ctx, &Msg{Title: rpcReqPrefix + name, Reply: rpcResPrefix + reply, Body: data}); err != nil {
		return nil, err
	}

	// Pick real deadline
	after := time.Second // TODO: make larger
	deadline, ok := ctx.Deadline()
	if ok {
		after = deadline.Sub(time.Now())
	}

	// Whichever comes first
	select {
	case bits := <-resc:
		return bits, nil
	case <-time.After(after):
		return nil, errors.New("Request Timeout")
	}
}

// Handle provides a response to a fn
func (s *Conn) Handle(name string, fn func(Context, []byte) ([]byte, error)) (Stream, error) {
	if s.err != nil || s.trans == nil {
		return nil, s.err
	}
	return s.Subscribe(rpcReqPrefix+name, func(data []byte) {
		id := string(data[:8])
		res, err := fn(nil, data[8:]) // TODO: pass in context
		if err != nil {
			res = append([]byte{0}, []byte(err.Error())...)
		} else {
			res = append([]byte{1}, res...)
		}
		s.Publish(nil, rpcResPrefix+id, res) // TODO: handle error here
	})
}

// Subscribe to a given broadcast channel
func (s *Conn) Subscribe(name string, cb func([]byte)) (Stream, error) {
	if s.err != nil || s.trans == nil {
		return nil, s.err
	}
	s.hands[name] = cb
	return func() error {
		delete(s.hands, name)
		return nil
	}, nil
}

// Publish some data!
func (s *Conn) Publish(ctx Context, name string, data []byte) error {
	if s.err != nil || s.trans == nil {
		return s.err
	}
	// checkPubSubName(name)
	// v := len(name) // Thanks binary.LittleEndian
	// message := append([]byte{byte(v), byte(v >> 8)}, data...)
	return s.trans.Push(ctx, &Msg{Title: name, Body: data})
}

// Close kills a connection
func (s *Conn) Close() error {
	killer := make(chan error, 1)
	s.close <- killer
	return <-killer // TODO: timeout
}

// func checkPubSubName(name string) {
// 	if len(name) > 16 {
// 		panic("Publish name too long: '" + name + "' > 16")
// 	}
// }
