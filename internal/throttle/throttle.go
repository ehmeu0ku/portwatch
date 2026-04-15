// Package throttle provides per-key event throttling to suppress
// repeated alerts for the same port within a configurable window.
package throttle

import (
	"sync"
	"time"
)

// Throttle tracks the last emission time for a string key and suppresses
// subsequent calls until the cooldown window has elapsed.
type Throttle struct {
	mu       sync.Mutex
	last     map[string]time.Time
	window   time.Duration
	nowFn    func() time.Time
}

// New returns a Throttle with the given suppression window.
func New(window time.Duration) *Throttle {
	return &Throttle{
		last:   make(map[string]time.Time),
		window: window,
		nowFn:  time.Now,
	}
}

// Allow returns true if the key has not been seen within the current window,
// and records the current time for the key. Returns false if the key is
// still within its cooldown period.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.nowFn()
	if last, ok := t.last[key]; ok {
		if now.Sub(last) < t.window {
			return false
		}
	}
	t.last[key] = now
	return true
}

// Reset removes the throttle record for key, allowing the next call to
// Allow for that key to pass immediately.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, key)
}

// Len returns the number of keys currently tracked.
func (t *Throttle) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.last)
}

// Purge removes all keys whose last-seen time is older than the window,
// freeing memory for ports that are no longer active.
func (t *Throttle) Purge() {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.nowFn()
	for k, last := range t.last {
		if now.Sub(last) >= t.window {
			delete(t.last, k)
		}
	}
}
