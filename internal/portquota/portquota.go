// Package portquota enforces per-protocol port count quotas and emits
// an alert when the number of observed listeners on a given protocol
// exceeds a configured ceiling.
package portquota

import (
	"fmt"
	"sync"
)

// Entry holds the current count and ceiling for a single protocol key.
type Entry struct {
	Count   int
	Ceiling int
}

// Quota tracks how many ports are active per protocol and whether any
// protocol has breached its ceiling.
type Quota struct {
	mu      sync.Mutex
	counts  map[string]int
	ceilings map[string]int
	default_ int
}

// New creates a Quota with the given default ceiling applied to any
// protocol that does not have an explicit ceiling registered.
func New(defaultCeiling int) *Quota {
	return &Quota{
		counts:   make(map[string]int),
		ceilings: make(map[string]int),
		default_: defaultCeiling,
	}
}

// SetCeiling registers an explicit ceiling for a protocol (e.g. "tcp", "udp").
func (q *Quota) SetCeiling(proto string, ceiling int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.ceilings[proto] = ceiling
}

// Observe increments the active count for proto and reports whether the
// ceiling has been breached. It returns true only on the transition from
// at-or-below to above the ceiling (i.e. exactly when count > ceiling).
func (q *Quota) Observe(proto string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.counts[proto]++
	ceiling := q.ceiling(proto)
	return q.counts[proto] > ceiling
}

// Release decrements the active count for proto. It is a no-op if the
// count is already zero.
func (q *Quota) Release(proto string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.counts[proto] > 0 {
		q.counts[proto]--
	}
}

// Exceeds reports whether the current count for proto is above its ceiling.
func (q *Quota) Exceeds(proto string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.counts[proto] > q.ceiling(proto)
}

// Snapshot returns a copy of all tracked entries.
func (q *Quota) Snapshot() map[string]Entry {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make(map[string]Entry, len(q.counts))
	for proto, count := range q.counts {
		out[proto] = Entry{Count: count, Ceiling: q.ceiling(proto)}
	}
	return out
}

// String returns a human-readable summary of all quota entries.
func (q *Quota) String() string {
	snap := q.Snapshot()
	s := ""
	for proto, e := range snap {
		s += fmt.Sprintf("%s: %d/%d ", proto, e.Count, e.Ceiling)
	}
	return s
}

// ceiling returns the ceiling for proto, falling back to the default.
// Must be called with q.mu held.
func (q *Quota) ceiling(proto string) int {
	if c, ok := q.ceilings[proto]; ok {
		return c
	}
	return q.default_
}
