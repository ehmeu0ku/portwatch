// Package portexpiry tracks how long a port has been continuously absent
// and emits an expiry signal once the absence exceeds a configured TTL.
package portexpiry

import (
	"sync"
	"time"
)

// Entry records when a port first went absent.
type Entry struct {
	FirstAbsent time.Time
}

// Tracker monitors absent ports and reports when they have expired.
type Tracker struct {
	mu  sync.Mutex
	ttl time.Duration
	now func() time.Time
	map_ map[string]Entry
}

// New returns a Tracker with the given absence TTL.
func New(ttl time.Duration) *Tracker {
	return &Tracker{
		ttl:  ttl,
		now:  time.Now,
		map_: make(map[string]Entry),
	}
}

// Observe records that key is absent. Returns true once the key has been
// absent for longer than the configured TTL.
func (t *Tracker) Observe(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	e, ok := t.map_[key]
	if !ok {
		t.map_[key] = Entry{FirstAbsent: now}
		return false
	}
	return now.Sub(e.FirstAbsent) >= t.ttl
}

// Forget removes the key from tracking (e.g. port came back).
func (t *Tracker) Forget(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.map_, key)
}

// Len returns the number of currently tracked absent ports.
func (t *Tracker) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.map_)
}
