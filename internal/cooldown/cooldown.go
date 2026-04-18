// Package cooldown tracks per-key cooldown windows, preventing repeated
// actions from firing more often than a configured duration allows.
package cooldown

import (
	"sync"
	"time"
)

// Cooldown tracks the last fire time for arbitrary string keys.
type Cooldown struct {
	mu       sync.Mutex
	last     map[string]time.Time
	window   time.Duration
	nowFunc  func() time.Time
}

// New returns a Cooldown with the given window duration.
func New(window time.Duration) *Cooldown {
	return &Cooldown{
		last:    make(map[string]time.Time),
		window:  window,
		nowFunc: time.Now,
	}
}

// Allow returns true if the key has not fired within the cooldown window.
// If allowed, it records the current time as the last fire time.
func (c *Cooldown) Allow(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.nowFunc()
	if t, ok := c.last[key]; ok && now.Sub(t) < c.window {
		return false
	}
	c.last[key] = now
	return true
}

// Reset clears the recorded fire time for the given key, allowing it
// to fire immediately on the next call to Allow.
func (c *Cooldown) Reset(key string) {
	c.mu.Lock()
	delete(c.last, key)
	c.mu.Unlock()
}

// Len returns the number of keys currently tracked.
func (c *Cooldown) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.last)
}
