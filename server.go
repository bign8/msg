package msg

import "errors"

// Locker is a copy of sync.Locker.
type Locker interface {
	Lock()
	Unlock()
}

// NewServer constructs a new service with a given locker.
func NewServer(read, write Locker) *Server {
	if read == nil || write == nil {
		panic("NewServer: invalid lockers provided")
	}
	return &Server{
		handlers: make(map[string]*service),
		r:        read,
		w:        write,
	}
}

// Server is designed to handle call requests.
type Server struct {
	handlers map[string]*service
	r, w     Locker
}

type service struct {
	sd ServiceDesc
	ss interface{}
}

// Register allows a server to handle requests for a given service.
func (s *Server) Register(sd ServiceDesc, ss interface{}) {
	s.w.Lock()
	s.handlers[sd.Name] = &service{sd: sd, ss: ss}
	s.w.Unlock()
}

// Call invokes a function on a server
func (s *Server) Call(ctx Context, service, method string, in Reader, out Writer) error {
	s.r.Lock()
	svc, ok := s.handlers[service]
	s.r.Unlock()
	if !ok {
		return errors.New("server: unregistered service: " + service)
	}
	var target MethodHandler
	for _, m := range svc.sd.Methods {
		if m.Name == method {
			target = m.Handler
			break
		}
	}
	if target == nil {
		return errors.New("server: invalid method: '" + method + "' on service '" + service + "'")
	}
	return target(ctx, svc.ss, in, out)
}

// ServiceDesc defines the adapter part of a service that allows Call to invoke methods on a service object.
type ServiceDesc struct {
	Name    string
	Methods []MethodDesc
}

// MethodDesc defines how a method can be invoked for a given
type MethodDesc struct {
	Name    string
	Handler MethodHandler
}

// MethodHandler is the core method handler for a given type
type MethodHandler func(ctx Context, svc interface{}, in Reader, out Writer) error
