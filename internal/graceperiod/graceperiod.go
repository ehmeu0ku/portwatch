// Package graceperiod suppresses alerts for ports that appear and disappear
// within a short observation window, reducing noise from transient listeners.
package graceperiod

import (
	"sync"
	"time"
)

// Entry tracks when a key was first seen.
type entry struct {
	firstSeen time.Time
}

// Filter holds pending keys and releases them only after the grace window
// has elapsed.
type Filter struct {
	mu      sync.Mutex
	window  time.Duration
	pending map[string]entry
	now     func() time.Time
}

// New creates a Filter with the given grace window.
func New(window time.Duration) *Filter {
	return &Filter{
		window:  window,
		pending: make(map[string]entry),
		now:     time.Now,
	}
}

// Observe registers key as seen. Returns true when the key has been present
// for longer than the grace window and should be forwarded.
func (f *Filter) Observe(key string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := f.now()
	e, ok := f.pending[key]
	if !ok {
		f.pending[key] = entry{firstSeen: now}
		return false
	}
	return now.Sub(e.firstSeen) >= f.window
}

// Forget removes key from the pending set (e.g. when a port disappears).
func (f *Filter) Forget(key string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.pending, key)
}

// Len returns the number of keys currently in the pending set.
func (f *Filter) Len() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.pending)
}
