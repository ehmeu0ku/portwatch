// Package portlifecycle tracks the full lifecycle of observed ports,
// recording first-seen, last-seen, and total observation count.
package portlifecycle

import (
	"fmt"
	"sync"
	"time"
)

// Entry holds lifecycle metadata for a single port key.
type Entry struct {
	FirstSeen  time.Time
	LastSeen   time.Time
	SeenCount  int
	GoneAt     *time.Time
}

// Tracker maintains lifecycle state for each port key.
type Tracker struct {
	mu      sync.Mutex
	entries map[string]*Entry
	now     func() time.Time
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{
		entries: make(map[string]*Entry),
		now:     time.Now,
	}
}

// Observe records an active observation for key.
// Returns the entry after updating it.
func (t *Tracker) Observe(key string) *Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	e, ok := t.entries[key]
	if !ok {
		e = &Entry{FirstSeen: now}
		t.entries[key] = e
	}
	e.LastSeen = now
	e.SeenCount++
	e.GoneAt = nil
	return copyEntry(e)
}

// MarkGone records that the port is no longer observed.
func (t *Tracker) MarkGone(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[key]
	if !ok {
		return
	}
	now := t.now()
	e.GoneAt = &now
}

// Get returns the entry for key and whether it exists.
func (t *Tracker) Get(key string) (*Entry, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[key]
	if !ok {
		return nil, false
	}
	return copyEntry(e), true
}

// Forget removes the entry for key.
func (t *Tracker) Forget(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, key)
}

// Len returns the number of tracked keys.
func (t *Tracker) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.entries)
}

// String returns a human-readable summary of the entry.
func (e *Entry) String() string {
	if e.GoneAt != nil {
		return fmt.Sprintf("seen=%d first=%s last=%s gone=%s",
			e.SeenCount, e.FirstSeen.Format(time.RFC3339),
			e.LastSeen.Format(time.RFC3339), e.GoneAt.Format(time.RFC3339))
	}
	return fmt.Sprintf("seen=%d first=%s last=%s",
		e.SeenCount, e.FirstSeen.Format(time.RFC3339),
		e.LastSeen.Format(time.RFC3339))
}

func copyEntry(e *Entry) *Entry {
	copy := *e
	if e.GoneAt != nil {
		t := *e.GoneAt
		copy.GoneAt = &t
	}
	return &copy
}
