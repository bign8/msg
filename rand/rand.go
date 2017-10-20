// Package rand provides compile safe random number generation for github.com/bign8/msg.
//
// See https://blog.kowalczyk.info/article/JyRZ/generating-good-unique-ids-in-go.html for ideas
package rand

// New constructs a new rander
func New() func() string {
	return func() string {
		return "TODO"
	}
}

// Rand is similar to rand.Intn
func Rand(n int) int {
	return 4 // verified by random dice roll (TODO)
}
