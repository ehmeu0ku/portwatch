// Package portage tracks how long a port has been continuously observed.
package portage

import (
	"sync"
	"time"
)

// Entry holds the first-seen timestamp and last-seen timestamp for a port key.
type Entry struct {
	FirstSeen time.Time
	LastSeen  time.Time
}

// Age returns how long the port has been continuously observed.
func (e Entry) Age(now time.Time) time.Duration {
	return now.Sub(e.FirstSeen)
}

// Tracker records when ports were first and last seen.
type Tracker struct {
	mu      sync.Mutex
	entries map[string]Entry
	now     func() time.Time
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{
		entries: make(map[string]Entry),
		now:     time.Now,
	}
}

// Observe records an observation for key. Returns the entry after update.
func (t *Tracker) Observe(key string) Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	e, ok := t.entries[key]
	if !ok {
		e = Entry{FirstSeen: now}
	}
	e.LastSeen = now
	t.entries[key] = e
	return e
}

// Forget removes the entry for key.
func (t *Tracker) Forget(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, key)
}

// Get returns the entry for key and whether it exists.
func (t *Tracker) Get(key string) (Entry, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[key]
	return e, ok
}

// Len returns the number of tracked ports.
func (t *Tracker) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.entries)
}
