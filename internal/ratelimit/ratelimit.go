// Package ratelimit provides a simple per-key rate limiter to suppress
// repeated alerts for the same port within a configurable cooldown window.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks the last alert time for each key and suppresses duplicates
// that occur within the cooldown duration.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	lastSeen map[string]time.Time
}

// New creates a Limiter with the given cooldown duration.
// Keys that were alerted within the cooldown window will be suppressed.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		lastSeen: make(map[string]time.Time),
	}
}

// Allow returns true if the key has not been seen within the cooldown window,
// and records the current time for that key. Returns false if the key is
// still within its cooldown period.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if last, ok := l.lastSeen[key]; ok {
		if now.Sub(last) < l.cooldown {
			return false
		}
	}
	l.lastSeen[key] = now
	return true
}

// Reset removes the cooldown record for the given key, allowing the next
// call to Allow for that key to pass immediately.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.lastSeen, key)
}

// Purge removes all keys whose last-seen time is older than the cooldown,
// keeping memory usage bounded during long-running daemon operation.
func (l *Limiter) Purge() {
	l.mu.Lock()
	defer l.mu.Unlock()

	cutoff := time.Now().Add(-l.cooldown)
	for key, t := range l.lastSeen {
		if t.Before(cutoff) {
			delete(l.lastSeen, key)
		}
	}
}
