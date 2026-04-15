// Package suppress provides a mechanism to suppress repeated alerts
// for the same port event within a configurable time window.
package suppress

import (
	"sync"
	"time"
)

// Key uniquely identifies a suppressible event.
type Key struct {
	Proto string
	Addr  string
	Port  uint16
	Kind  string // "new" or "gone"
}

// entry holds the expiry time for a suppressed key.
type entry struct {
	expiry time.Time
}

// Suppressor tracks which events have been recently seen and
// suppresses duplicates until their window expires.
type Suppressor struct {
	mu     sync.Mutex
	window time.Duration
	seen   map[Key]entry
}

// New creates a Suppressor with the given suppression window.
// Events with the same key will be suppressed for the duration
// of the window after they are first seen.
func New(window time.Duration) *Suppressor {
	return &Suppressor{
		window: window,
		seen:   make(map[Key]entry),
	}
}

// Allow returns true if the event identified by key should be
// forwarded, and false if it is currently suppressed.
// The first call for a key always returns true and starts the window.
func (s *Suppressor) Allow(k Key) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	if e, ok := s.seen[k]; ok && now.Before(e.expiry) {
		return false
	}
	s.seen[k] = entry{expiry: now.Add(s.window)}
	return true
}

// Reset removes a key from the suppression table, allowing the
// next event for that key to pass through immediately.
func (s *Suppressor) Reset(k Key) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.seen, k)
}

// Purge removes all expired entries from the suppression table.
// Call periodically to prevent unbounded memory growth.
func (s *Suppressor) Purge() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for k, e := range s.seen {
		if now.After(e.expiry) {
			delete(s.seen, k)
		}
	}
}

// Len returns the number of keys currently tracked (including expired).
func (s *Suppressor) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.seen)
}
