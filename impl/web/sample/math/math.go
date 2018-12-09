package math

import "github.com/bign8/msg"

// Math is the service interface for a math service.
type Math interface {

	// Round returns the nearest integer, rounding half away from zero.
	Round(msg.Context, []float64) ([]float64, error)
}
