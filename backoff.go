package backoff

import (
	"math/rand"
	"sync"
	"time"
)

// Backoff is a type that returns the next duration to wait before retrying
// or proceeding with an operation.
type Backoff interface {
	Next() time.Duration
}

// ConstantBackoff returns a Backoff that always returns the same duration.
func ConstantBackoff(dur time.Duration) Backoff {
	return constantBackoff{dur}
}

type constantBackoff struct {
	dur time.Duration
}

func (c constantBackoff) Next() time.Duration {
	return c.dur
}

// ExponentialBackoff returns a Backoff that exponentially increases the duration
// with a 25% jitter with each call to Next, up to a maximum duration.
//
// If the maxDelay is less than the initialDelay, ExponentialBackoff panics.
func ExponentialBackoff(initialDelay, maxDelay time.Duration) Backoff {
	if maxDelay < initialDelay {
		panic("maxDelay must be greater than initialDelay")
	}
	return &exponentialBackoff{
		dur:    initialDelay,
		max:    maxDelay,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
		mu:     sync.Mutex{},
	}
}

type exponentialBackoff struct {
	dur    time.Duration
	max    time.Duration
	random *rand.Rand
	mu     sync.Mutex
}

func (e *exponentialBackoff) Next() time.Duration {
	next := e.dur

	maxJitter := e.dur / 4
	e.mu.Lock()
	randomFactor := time.Duration(e.random.Intn(int(maxJitter*2+1))) - maxJitter
	e.mu.Unlock()
	e.dur = e.dur*2 + randomFactor
	if e.dur > e.max {
		e.dur = e.max
	}
	return next
}

// RandomBackoff returns a Backoff that returns a random duration between min
// and max with each call to Next.
//
// If max is less than min, RandomBackoff panics.
func RandomBackoff(min, max time.Duration) Backoff {
	if max < min {
		panic("max must be greater than min")
	}
	return &randomBackoff{
		min:  min,
		max:  max,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
		mu:   sync.Mutex{},
	}
}

type randomBackoff struct {
	min  time.Duration
	max  time.Duration
	rand *rand.Rand
	mu   sync.Mutex
}

func (r *randomBackoff) Next() time.Duration {
	r.mu.Lock()
	defer r.mu.Unlock()

	durRange := r.max - r.min
	return r.min + time.Duration(r.rand.Int63n(int64(durRange)))
}
