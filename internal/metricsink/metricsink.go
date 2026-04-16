// Package metricsink provides a lightweight in-process counter store that
// accumulates named uint64 metrics and exposes a snapshot for reporting.
package metricsink

import (
	"sync"
	"sync/atomic"
)

// Sink holds named counters.
type Sink struct {
	mu       sync.RWMutex
	counters map[string]*atomic.Uint64
}

// New returns an empty Sink.
func New() *Sink {
	return &Sink{counters: make(map[string]*atomic.Uint64)}
}

// Inc increments the named counter by 1.
func (s *Sink) Inc(name string) {
	s.counter(name).Add(1)
}

// Add adds delta to the named counter.
func (s *Sink) Add(name string, delta uint64) {
	s.counter(name).Add(delta)
}

// Get returns the current value of the named counter.
func (s *Sink) Get(name string) uint64 {
	s.mu.RLock()
	c, ok := s.counters[name]
	s.mu.RUnlock()
	if !ok {
		return 0
	}
	return c.Load()
}

// Snapshot returns a copy of all counters at this instant.
func (s *Sink) Snapshot() map[string]uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]uint64, len(s.counters))
	for k, v := range s.counters {
		out[k] = v.Load()
	}
	return out
}

// Reset sets the named counter back to zero.
func (s *Sink) Reset(name string) {
	s.counter(name).Store(0)
}

func (s *Sink) counter(name string) *atomic.Uint64 {
	s.mu.RLock()
	c, ok := s.counters[name]
	s.mu.RUnlock()
	if ok {
		return c
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if c, ok = s.counters[name]; ok {
		return c
	}
	c = &atomic.Uint64{}
	s.counters[name] = c
	return c
}
