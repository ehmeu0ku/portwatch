// Package scorecard tracks per-port event counts and computes a simple
// risk score that downstream components can use for prioritisation.
package scorecard

import (
	"sync"
	"time"
)

// Entry holds the running totals for a single port key.
type Entry struct {
	Key        string
	SeenCount  int
	AlertCount int
	LastSeen   time.Time
	Score      float64
}

// Scorecard maintains a thread-safe map of port keys to entries.
type Scorecard struct {
	mu      sync.RWMutex
	entries map[string]*Entry
}

// New returns an initialised Scorecard.
func New() *Scorecard {
	return &Scorecard{entries: make(map[string]*Entry)}
}

// Record increments the seen counter for key and optionally the alert
// counter when alerted is true, then recalculates the score.
func (s *Scorecard) Record(key string, alerted bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.entries[key]
	if !ok {
		e = &Entry{Key: key}
		s.entries[key] = e
	}
	e.SeenCount++
	e.LastSeen = time.Now()
	if alerted {
		e.AlertCount++
	}
	e.Score = score(e)
}

// Get returns the entry for key and whether it exists.
func (s *Scorecard) Get(key string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[key]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// Snapshot returns a copy of all current entries.
func (s *Scorecard) Snapshot() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, *e)
	}
	return out
}

// Reset removes the entry for key.
func (s *Scorecard) Reset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, key)
}

// score computes a simple 0-100 risk score based on alert ratio.
func score(e *Entry) float64 {
	if e.SeenCount == 0 {
		return 0
	}
	ratio := float64(e.AlertCount) / float64(e.SeenCount)
	return ratio * 100
}
