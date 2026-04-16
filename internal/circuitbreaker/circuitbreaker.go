// Package circuitbreaker implements a simple circuit breaker that opens
// after a threshold of consecutive failures and resets after a cooldown.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned when the circuit is open and calls are rejected.
var ErrOpen = errors.New("circuit breaker is open")

// State represents the current circuit breaker state.
type State int

const (
	StateClosed State = iota
	StateOpen
)

// Breaker tracks consecutive failures and opens the circuit when the
// failure threshold is exceeded.
type Breaker struct {
	mu        sync.Mutex
	threshold int
	cooldown  time.Duration
	failures  int
	openedAt  time.Time
	state     State
}

// New returns a Breaker that opens after threshold consecutive failures
// and resets after cooldown.
func New(threshold int, cooldown time.Duration) *Breaker {
	return &Breaker{
		threshold: threshold,
		cooldown:  cooldown,
	}
}

// Allow returns nil if the call should proceed, or ErrOpen if the circuit
// is open. It automatically transitions back to closed after the cooldown.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.state == StateOpen {
		if time.Since(b.openedAt) >= b.cooldown {
			b.state = StateClosed
			b.failures = 0
		} else {
			return ErrOpen
		}
	}
	return nil
}

// RecordSuccess resets the failure counter.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
}

// RecordFailure increments the failure counter and opens the circuit if
// the threshold is reached.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.threshold {
		b.state = StateOpen
		b.openedAt = time.Now()
	}
}

// State returns the current state of the breaker.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

// Reset forces the breaker back to closed with zero failures.
func (b *Breaker) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.state = StateClosed
	b.failures = 0
}
