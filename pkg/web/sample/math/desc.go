package math

import "github.com/bign8/msg"

var MathDesc = msg.ServiceDesc{
	Name: "math",
	Methods: []msg.MethodDesc{{
		Name: "round",
		Args: msg.Slice{
			Value: msg.Float64,
		},
		Result: msg.Slice{
			Value: msg.Float64,
		},
		Handler: func(ctx msg.Context, svc interface{}, args interface{}) (interface{}, error) {
			return svc.(Math).Round(ctx, args.([]float64))
		},
	}},
}
