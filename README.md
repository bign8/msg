# MSG
A generic abstraction over various methods of communication.

## Encoding
At some point objects need to get serialized to bytes in preporation for transport.

```go
type Encoding interface {
    Encode(io.Writer, interface{}) error
    Decode(io.Reader, interface{}) error
}
```

## Transport
Once data is converted to a data stream, we have to send that stream somewhere.
Transports are responsible for discovering downstream services and can be used to exchange supported encodings for a service.

```go
// type Handler func(context.Context, io.Reader) error

type Transport interface {
    // Send(ctx context.Context, service, method string, msg io.Reader) error
    // Recv(ctx context.Context, service, method string, cb Handler) error
    Call(ctx context.Context, args io.Reader, res io.Writer) error
}
```

## Discovery

```go
type Discovery interface {
    Connect(ctx context.Context, service, method string) (Transport, error)
}
```

## Usage
To make things nice for the consumer, but composable for designers, the usage abstracts away `Transport` and `Encoding` to a core function that can be called from anywhere.  This function uses a servers `context.Context` which will have both a `Transport` and `Encoding` configured and can allow simple calls to other services.

```go
func Call(ctx context.Context, service, method string, args, res interface{}) error {
    // TODO
}

// Configuring a service context
func main() {
    caller := msg.NewCaller(json.Encoding, http.REST)

    ctx := context.TODO() // can be provided from r.Context() or otherwise
    ctx = msg.WithCaller(ctx, caller) // can be bound via the default-serve-mux

    req := &MyRequest{
        Message: `hello-world`,
    }
    res := &MyResponse{}
    err := Call(ctx, `example`, `hello-world`, req, res)
    fmt.Printf("Got Response: %v %v", res, err)
}
```