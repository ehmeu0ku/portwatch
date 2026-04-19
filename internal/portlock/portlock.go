// Package portlock tracks which ports are "locked" (expected to be
// permanently open) and which are transient, helping reduce noise for
// well-known long-running services.
package portlock

import (
	"sync"
	"time"
)

// Entry records when a port was first seen and whether it has been
// promoted to locked status.
type Entry struct {
	FirstSeen time.Time
	LockedAt  time.Time
	Locked    bool
}

// Store manages port lock state.
type Store struct {
	mu       sync.Mutex
	entries  map[string]*Entry
	minAge   time.Duration
}

// New returns a Store that promotes a port to locked after minAge of
// continuous presence.
func New(minAge time.Duration) *Store {
	return &Store{
		entries: make(map[string]*Entry),
		minAge:  minAge,
	}
}

// Observe records a sighting of key (e.g. "tcp:443"). It returns true
// if the port transitions to locked on this call.
func (s *Store) Observe(key string, now time.Time) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.entries[key]
	if !ok {
		s.entries[key] = &Entry{FirstSeen: now}
		return false
	}
	if e.Locked {
		return false
	}
	if now.Sub(e.FirstSeen) >= s.minAge {
		e.Locked = true
		e.LockedAt = now
		return true
	}
	return false
}

// IsLocked reports whether key is currently locked.
func (s *Store) IsLocked(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[key]
	return ok && e.Locked
}

// Release removes key from the store (port gone).
func (s *Store) Release(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, key)
}

// Snapshot returns a copy of all entries.
func (s *Store) Snapshot() map[string]Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[string]Entry, len(s.entries))
	for k, v := range s.entries {
		out[k] = *v
	}
	return out
}
