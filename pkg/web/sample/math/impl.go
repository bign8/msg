package math

import (
	"math"

	"github.com/bign8/msg"
)

type mathService struct{}

func (m *mathService) Round(ctx msg.Context, in []float64) ([]float64, error) {
	out := make([]float64, len(in))
	for i, v := range in {
		out[i] = math.Round(v)
	}
	return out, nil
}
