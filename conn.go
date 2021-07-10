package msg

import (
	"errors"
	"time"
)

// New builds a new Sock that performs automatic retires on connection failures
func New(trans Transport) *Conn {
	conn := &Conn{
		hands: make(map[string]Handler, 1),
		trans: trans,
		close: make(chan chan error, 1),
		err:   ErrClosed,
	}
	return conn
}

// Conn is a core managed connection for all communication
type Conn struct {
	hands map[string]Handler // array of active handlers
	trans Transport          // currently active transport

	// TODO: think about removing these guys
	err   error
	att   int             // retry attempt
	close chan chan error // closing
}

// Open starts a connection
func (s *Conn) Open() error {
	s.err = s.trans.Open()
	if s.err != nil {
		return s.err
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
	fn(msg)
}

func (s *Conn) genID() string {
	return "01234567" // NOTE: should always be len(8) so decoding is consistent
}

// Request executes an RPC
func (s *Conn) Request(ctx ContextOld, msg *Msg) (*Msg, error) {
	if s.err != nil || s.trans == nil {
		return nil, s.err
	}
	if !s.trans.Able() {
		return nil, ErrClosed
	}
	resc := make(chan *Msg, 1)
	defer close(resc) // yay memory leaks

	// Generate response subscription
	reply := rpcResPrefix + s.genID()
	sub, err := s.Subscribe(reply, func(res *Msg) error {
		resc <- res
		return nil
	})
	if err != nil {
		return nil, err
	}
	defer sub.Close() // TODO: log error here
	msg.Reply = reply

	// Publish request to the cloud
	msg.Title = rpcReqPrefix + msg.Title
	if err := s.trans.Push(ctx, msg); err != nil {
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
	case res := <-resc:
		return res, nil
	case <-time.After(after):
		return nil, errors.New("Request Timeout")
	}
}

// Handle provides a response to a fn
func (s *Conn) Handle(name string, fn func(ContextOld, *Msg) (*Msg, error)) (Stream, error) {
	if s.err != nil || s.trans == nil {
		return nil, s.err
	}
	return s.Subscribe(rpcReqPrefix+name, func(msg *Msg) error {
		res, err := fn(nil, msg) // TODO: pass in context
		if err != nil {
			return err
		}
		res.Title = rpcResPrefix + res.Title
		return s.trans.Push(nil, res) // TODO: handle error here
	})
}

// Subscribe to a given broadcast channel
func (s *Conn) Subscribe(name string, cb Handler) (Stream, error) {
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
func (s *Conn) Publish(ctx ContextOld, msg *Msg) error {
	if s.err != nil || s.trans == nil {
		return s.err
	}
	return s.trans.Push(ctx, msg)
}

// Close kills a connection
func (s *Conn) Close() error {
	killer := make(chan error, 1)
	s.close <- killer
	return <-killer // TODO: timeout
}
