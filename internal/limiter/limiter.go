// Package limiter provides a token-bucket style rate limiter that caps
// the number of events dispatched per key within a sliding window.
package limiter

import (
	"sync"
	"time"
)

// Limiter tracks per-key event counts and enforces a maximum burst
// within a rolling window duration.
type Limiter struct {
	mu      sync.Mutex
	window  time.Duration
	maxBurst int
	buckets map[string]*bucket
}

type bucket struct {
	count  int
	reset  time.Time
}

// New creates a Limiter that allows at most maxBurst events per key
// within the given window.
func New(window time.Duration, maxBurst int) *Limiter {
	return &Limiter{
		window:   window,
		maxBurst: maxBurst,
		buckets:  make(map[string]*bucket),
	}
}

// Allow returns true if the event for key is within the burst limit.
// It increments the counter and resets it when the window expires.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	b, ok := l.buckets[key]
	if !ok || now.After(b.reset) {
		l.buckets[key] = &bucket{count: 1, reset: now.Add(l.window)}
		return true
	}
	if b.count >= l.maxBurst {
		return false
	}
	b.count++
	return true
}

// Reset clears the bucket for key, allowing immediate passage on the
// next call to Allow.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.buckets, key)
}

// Len returns the number of active buckets.
func (l *Limiter) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.buckets)
}
