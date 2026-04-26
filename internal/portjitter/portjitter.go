// Package portjitter tracks timing variance between successive observations
// of a port. A port that appears and disappears erratically has high jitter;
// one that remains stable has low jitter. This can surface ephemeral or
// scanning behaviour that other detectors miss.
package portjitter

import (
	"sync"
	"time"
)

// Entry holds the rolling statistics for a single port key.
type Entry struct {
	lastSeen  time	.Time
	intervals []time.Duration
	cap       int
}

// Tracker measures inter-observation jitter for port keys.
type Tracker struct {
	mu      sync.Mutex
	entries map[string]*Entry
	window  int // max intervals to retain
}

// New returns a Tracker that retains at most window intervals per key.
// A window of zero is replaced with the default of 10.
func New(window int) *Tracker {
	if window <= 0 {
		window = 10
	}
	return &Tracker{
		entries: make(map[string]*Entry),
		window:  window,
	}
}

// Observe records a new observation for key at the given time.
// It returns the jitter (stddev-like spread) of recent intervals in
// nanoseconds, and a boolean indicating whether enough samples exist
// (at least 2 intervals) to produce a meaningful value.
func (t *Tracker) Observe(key string, now time.Time) (jitter time.Duration, stable bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	e, ok := t.entries[key]
	if !ok {
		t.entries[key] = &Entry{lastSeen: now, cap: t.window}
		return 0, false
	}

	if !e.lastSeen.IsZero() {
		interval := now.Sub(e.lastSeen)
		e.intervals = append(e.intervals, interval)
		if len(e.intervals) > e.cap {
			e.intervals = e.intervals[len(e.intervals)-e.cap:]
		}
	}
	e.lastSeen = now

	if len(e.intervals) < 2 {
		return 0, false
	}

	return spread(e.intervals), true
}

// Forget removes tracking state for key.
func (t *Tracker) Forget(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, key)
}

// Len returns the number of keys currently tracked.
func (t *Tracker) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.entries)
}

// spread returns the max deviation from the mean of a slice of durations.
func spread(intervals []time.Duration) time.Duration {
	var sum int64
	for _, d := range intervals {
		sum += int64(d)
	}
	mean := sum / int64(len(intervals))
	var maxDev int64
	for _, d := range intervals {
		dev := int64(d) - mean
		if dev < 0 {
			dev = -dev
		}
		if dev > maxDev {
			maxDev = dev
		}
	}
	return time.Duration(maxDev)
}
