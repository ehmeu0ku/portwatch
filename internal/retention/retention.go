// Package retention enforces a maximum age on audit and history entries,
// pruning records older than a configured TTL.
package retention

import (
	"sync"
	"time"
)

// Entry is any record that carries a timestamp.
type Entry interface {
	Timestamp() time.Time
}

// Store is a generic append-and-prune store for timestamped entries.
type Store struct {
	mu      sync.Mutex
	entries []Entry
	ttl     time.Duration
	now     func() time.Time
}

// New returns a Store that discards entries older than ttl.
func New(ttl time.Duration) *Store {
	return &Store{ttl: ttl, now: time.Now}
}

// Add appends e and immediately evicts entries beyond ttl.
func (s *Store) Add(e Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, e)
	s.prune()
}

// Entries returns a copy of all non-expired entries.
func (s *Store) Entries() []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prune()
	out := make([]Entry, len(s.entries))
	copy(out, s.entries)
	return out
}

// Len returns the number of live entries.
func (s *Store) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prune()
	return len(s.entries)
}

// prune must be called with s.mu held.
func (s *Store) prune() {
	cutoff := s.now().Add(-s.ttl)
	i := 0
	for i < len(s.entries) && s.entries[i].Timestamp().Before(cutoff) {
		i++
	}
	s.entries = s.entries[i:]
}
