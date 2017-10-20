package msg

import "github.com/bign8/msg/rand"

const (
	retryCap  = 60 * 1000 // one-minute max retry
	retryBase = 50        // 50ms base retrys
)

// retryDelay gives the current retry delay in milliseconds
// TODO: https://www.awsarchitectureblog.com/2015/03/backoff.html
// TODO: benchmark
func retryDelay(attempt int) (delay, next int) {
	delay = 1 << uint(attempt-1) // 2^retryAttp
	delay *= retryBase
	if delay > retryCap {
		delay = retryCap
	}
	return int(rand.Rand(delay)), attempt + 1
}

// Backoff performs connnection retries on a given transport
func Backoff(baseMS, capMS int) Option {
	return func(c *Conn) {
		// builder := c.build
	}
}

type backoff struct {
	Transport
	att int
}

func (b backoff) Recv(fn func(string, []byte)) error {
	b.att = 0
	err := b.Transport.Recv(fn)
	if err != ErrClosed {
		return err
	}
	return err
}
