package msg

import (
	"errors"
	"time"
)

// New builds a new Sock that performs automatic retires on connection failures
func New(addr string, builder Builder, opts ...Option) *Conn {
	conn := &Conn{
		hands: make(map[string]func([]byte), 1),
		build: builder,
		addr:  addr,
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
	build Builder                 // How to construct a new Connection
	addr  string                  // address of connection

	// TODO: think about removing these guys
	err   error
	att   int             // retry attempt
	close chan chan error // closing
}

// Open is a blocking call that opens a connection
func (s *Conn) Open() error {
	return s.start()
}

func (s *Conn) recv(subject string, data []byte) {
	s.att = 0
	s.err = nil
	fn, ok := s.hands[subject]
	if !ok {
		print("Unsupported function:" + subject)
		return
	}
	fn(data)
}

func (s *Conn) start() error {
	s.trans, s.err = s.build(s.addr)
	if s.err != nil {
		return s.err
	}

	err := s.trans.Recv(s.recv)
	if err != ErrClosed {
		return err
	}

	var delay int
	delay, s.att = retryDelay(s.att)
	time.Sleep(time.Duration(delay))
	return s.start()
}

func (s *Conn) genID() string {
	return "01234567" // NOTE: should always be len(8) so decoding is consistent
}

// Request executes an RPC
func (s *Conn) Request(ctx Context, name string, data []byte) ([]byte, error) {
	if s.err != nil || s.trans == nil {
		return nil, s.err
	}
	resc := make(chan []byte, 1)
	defer close(resc) // yay memory leaks

	// Generate response subscription
	reply := s.genID()
	sub, err := s.Subscribe(rpcResPrefix+reply, func(res []byte) { resc <- res })
	if err != nil {
		return nil, err
	}
	defer sub.Unsubscribe() // TODO: log error here

	// Add reply channel to the request
	message := append([]byte(reply), data...) // TODO: send timeout too

	// Publish request to the cloud
	if err := s.Publish(rpcReqPrefix+name, message); err != nil {
		return nil, err
	}

	// Pick real deadline
	after := time.Minute
	deadline, ok := ctx.Deadline()
	if ok {
		after = deadline.Sub(time.Now())
	}

	// Whichever comes first
	select {
	case bits := <-resc:
		if bits[0] == 1 {
			return bits[1:], nil
		}
		return nil, errors.New(string(bits[1:]))
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
		s.Publish(rpcResPrefix+id, res) // TODO: handle error here
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
func (s *Conn) Publish(name string, data []byte) error {
	if s.err != nil || s.trans == nil {
		return s.err
	}
	// checkPubSubName(name)
	// v := len(name) // Thanks binary.LittleEndian
	// message := append([]byte{byte(v), byte(v >> 8)}, data...)
	return s.trans.Send(name, data)
}

// func checkPubSubName(name string) {
// 	if len(name) > 16 {
// 		panic("Publish name too long: '" + name + "' > 16")
// 	}
// }
