// Package test defines many helpers for verifying client/server call interactions.
package test

import "github.com/bign8/msg/test/internal"

//go:generate msg SimpleService

// "github.com/bign8/msg"
// "github.com/bign8/msg/iso/buff"

// type registrator func(msg.ServiceDesc, interface{})
//
// // RegisterSimpleService adds a simple service handler to a server.
// func RegisterSimpleService(service string, registrar registrator) {
// 	// ...
// }
//
// func simplePingHandler(ctx msg.Context, svc interface{}, in msg.Reader, out msg.Writer) error {
// 	return svc.(SimpleService).Ping()
// }
//
// func simpleHelloHandler(ctx msg.Context, svc interface{}, in msg.Reader, out msg.Writer) error {
// 	inb := buff.NewReader(in)
// 	name := inb.ReadStr()
// 	if err := inb.Err(); err != nil {
// 		return err
// 	}
// 	msg, err := svc.(SimpleService).Hello(name)
// 	if err != nil {
// 		return err
// 	}
// 	ob := buff.NewWriter(out)
// 	ob.WriteStr(msg)
// 	return ob.Err()
// }
//
// func simpleRPCHandler(ctx msg.Context, svc interface{}, in msg.Reader, out msg.Writer) error {
// 	inb := buff.NewReader(in)
// 	req := &Request{}
// 	inb.ReadType(req)
// 	if err := inb.Err(); err != nil {
// 		return err
// 	}
// 	res, err := svc.(SimpleService).RPC(req)
// 	if err != nil {
// 		return err
// 	}
// 	ob := buff.NewWriter(out)
// 	ob.WriteType(res)
// 	return ob.Err()
// }

// SimpleService is a basic service.
type SimpleService interface {
	Ping() error
	Hello(name string) (msg string, err error)
	RPC(req *Request) (res *Response, err error)
	Eat(arg internal.Grapes) error
}

// Request ...
type Request struct {
	Service string
	Method  string
	Args    []byte
}

// func (r *Request) Read(buf buff.Buffer)  { /* TODO */ }
// func (r *Request) Write(buf buff.Buffer) { /* TODO */ }

// Response ...
type Response struct {
	Res []byte
	// TODO: more
}

// func (r *Response) Read(buf buff.Buffer)  { /* TODO */ }
// func (r *Response) Write(buf buff.Buffer) { /* TODO */ }

// // SimpleServiceDesc ...
// var SimpleServiceDesc = msg.ServiceDesc{
// 	Name: "simple",
// 	Methods: []msg.MethodDesc{{
// 		Name:    "ping",
// 		Handler: simplePingHandler,
// 	}, {
// 		Name:    "hello",
// 		Handler: simpleHelloHandler,
// 	}, {
// 		Name:    "rpc",
// 		Handler: simpleRPCHandler,
// 	}},
// }
