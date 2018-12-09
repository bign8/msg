package msg

import "testing"

type noopLock int

func (noopLock) Lock()   {}
func (noopLock) Unlock() {}

func Test(t *testing.T) {

	// Create a new server
	s := NewServer(noopLock(1), noopLock(2))

	// Register a service handler on that server
	s.Register(ServiceDesc{
		Name: "test",
		Methods: []MethodDesc{{
			Name: "ping",
			Handler: func(ctx Context, svc interface{}, in Reader, out Writer) error {
				return nil
			},
		}},
	}, nil)

	{ // Invoke an invalid Service
		err := s.Call(Background, "rand", "get", nil, nil)
		if err == nil || err.Error() != "server: unregistered service: rand" {
			t.Errorf("Unexpected Error: %v", err)
		}
	}

	{ // Invoke an invalid method
		err := s.Call(Background, "test", "pong", nil, nil)
		if err == nil || err.Error() != "server: invalid method: 'pong' on service 'test'" {
			t.Errorf("Unexpected Error: %v", err)
		}
	}

	{ // Invoke a method
		err := s.Call(Background, "test", "ping", nil, nil)
		if err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
	}

	{ // Perform call with unbound context
		err := Call(Background, "test", "ping", nil, nil)
		if err == nil || err.Error() != "call: no caller on context" {
			t.Errorf("Unexpected Error: %v", err)
		}
	}

	{ // Perform call with a bound context
		ctx := WithCaller(Background, s.Call)
		ctx = WithValue(ctx, "asdf", "jkl;")
		err := Call(ctx, "test", "ping", nil, nil)
		if err != nil {
			t.Errorf("Unexpected Error: %v", err)
		}
	}
}
