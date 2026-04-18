// Package backoff provides exponential backoff with jitter for retry logic.
package backoff

import (
	"math"
	"math/rand"
	"time"
)

// Config holds backoff parameters.
type Config struct {
	Initial    time.Duration
	Max        time.Duration
	Multiplier float64
	Jitter     float64 // fraction in [0, 1]
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Initial:    500 * time.Millisecond,
		Max:        30 * time.Second,
		Multiplier: 2.0,
		Jitter:     0.2,
	}
}

// Backoff tracks retry state for a single operation.
type Backoff struct {
	cfg     Config
	attempt int
}

// New creates a Backoff with the given config.
func New(cfg Config) *Backoff {
	return &Backoff{cfg: cfg}
}

// Next returns the duration to wait before the next retry and advances state.
func (b *Backoff) Next() time.Duration {
	base := float64(b.cfg.Initial) * math.Pow(b.cfg.Multiplier, float64(b.attempt))
	if base > float64(b.cfg.Max) {
		base = float64(b.cfg.Max)
	}
	jitter := base * b.cfg.Jitter * (rand.Float64()*2 - 1)
	d := time.Duration(base + jitter)
	if d < 0 {
		d = b.cfg.Initial
	}
	b.attempt++
	return d
}

// Reset clears retry state.
func (b *Backoff) Reset() {
	b.attempt = 0
}

// Attempt returns the current attempt count.
func (b *Backoff) Attempt() int {
	return b.attempt
}

// NextWithMax returns the duration to wait, capped at the provided maximum.
// This is useful for callers that need a tighter bound for a specific retry.
func (b *Backoff) NextWithMax(max time.Duration) time.Duration {
	d := b.Next()
	if d > max {
		return max
	}
	return d
}
