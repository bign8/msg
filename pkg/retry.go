package msg

import "github.com/bign8/msg/pkg/rand"

const (
	retryCap  = 60 * 1000 // one-minute max retry
	retryBase = 50        // 50ms base retrys
)

// RetryDelay gives the current retry delay in milliseconds
// TODO: https://www.awsarchitectureblog.com/2015/03/backoff.html
// TODO: benchmark
func RetryDelay(attempt int) (delay, next int) {
	delay = 1 << uint(attempt-1) // 2^retryAttp
	delay *= retryBase
	if delay > retryCap {
		delay = retryCap
	}
	return int(rand.Rand(delay)), attempt + 1
}

// Maintain keeps a transport open with specific delay and cap requirements
func Maintain(t Transport, base, cap int) Transport {
	return &maintained{
		Transport: t,
		base:      base,
		cap:       cap,
	}
}

type maintained struct {
	Transport

	// retry attempt, base delay (raised to a power each retry), max delay
	attempt, base, cap int
}
