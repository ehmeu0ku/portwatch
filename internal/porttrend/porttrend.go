// Package porttrend tracks how frequently each port appears across scans
// and exposes a simple rising/falling trend signal.
package porttrend

import (
	"sync"
	"time"
)

// Direction indicates whether a port's presence is increasing or decreasing.
type Direction int

const (
	Stable  Direction = iota
	Rising            // seen more often recently
	Falling           // seen less often recently
)

func (d Direction) String() string {
	switch d {
	case Rising:
		return "rising"
	case Falling:
		return "falling"
	default:
		return "stable"
	}
}

// Sample is a single observation.
type Sample struct {
	At    time.Time
	Seen  bool
}

// Tracker records per-port samples and derives a trend.
type Tracker struct {
	mu      sync.Mutex
	window  int
	samples map[string][]Sample
}

// New creates a Tracker that keeps the last window samples per port key.
func New(window int) *Tracker {
	if window < 2 {
		window = 2
	}
	return &Tracker{window: window, samples: make(map[string][]Sample)}
}

// Record adds a sample for the given key.
func (t *Tracker) Record(key string, seen bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	s := t.samples[key]
	s = append(s, Sample{At: time.Now(), Seen: seen})
	if len(s) > t.window {
		s = s[len(s)-t.window:]
	}
	t.samples[key] = s
}

// Trend returns the current direction for the given key.
func (t *Tracker) Trend(key string) Direction {
	t.mu.Lock()
	defer t.mu.Unlock()
	s := t.samples[key]
	if len(s) < 2 {
		return Stable
	}
	half := len(s) / 2
	early := countSeen(s[:half])
	late := countSeen(s[half:])
	switch {
	case late > early:
		return Rising
	case late < early:
		return Falling
	default:
		return Stable
	}
}

func countSeen(ss []Sample) int {
	n := 0
	for _, s := range ss {
		if s.Seen {
			n++
		}
	}
	return n
}
