package main

import (
	"errors"

	"github.com/bign8/msg"
	"github.com/bign8/msg/pkg/http"
)

func main() {
	r := msg.Router(map[msg.ServiceMethod]msg.Callable{
		{`math`, `add`}: mathAdd,
	})
	http.Main(r)
}

func mathAdd(ctx msg.Context, service, method string, in msg.Reader, out msg.Writer) error {
	req := &MathAddRequest{}
	if err := req.ReadFrom(in); err != nil {
		return err
	}
	res := &MathAddResult{Result: req.Number1 + req.Number2}
	return res.WriteTo(out)
}

type MathAddRequest struct {
	Number1, Number2 int
}

func (req *MathAddRequest) ReadFrom(in msg.Reader) error {
	return errors.New(`TODO`)
}

type MathAddResult struct {
	Result int
}

func (res MathAddResult) WriteTo(out msg.Writer) error {
	return errors.New(`TODO`)
}
