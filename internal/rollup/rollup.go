// Package rollup groups rapid bursts of port events into a single
// summarised notification, reducing alert fatigue during restart storms.
package rollup

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/correlator"
)

// Group collects events that share the same port+kind key within a
// sliding window and flushes them as a slice when the window closes.
type Group struct {
	mu      sync.Mutex
	window  time.Duration
	buckets map[string][]correlator.Event
	timers  map[string]*time.Timer
	flush   func([]correlator.Event)
}

// New returns a Group that waits window duration of inactivity before
// calling flush with the accumulated events for that key.
func New(window time.Duration, flush func([]correlator.Event)) *Group {
	return &Group{
		window:  window,
		buckets: make(map[string][]correlator.Event),
		timers:  make(map[string]*time.Timer),
		flush:   flush,
	}
}

// Add inserts an event into its bucket and resets the flush timer.
func (g *Group) Add(ev correlator.Event) {
	key := eventKey(ev)

	g.mu.Lock()
	defer g.mu.Unlock()

	g.buckets[key] = append(g.buckets[key], ev)

	if t, ok := g.timers[key]; ok {
		t.Reset(g.window)
		return
	}

	g.timers[key] = time.AfterFunc(g.window, func() {
		g.mu.Lock()
		events := g.buckets[key]
		delete(g.buckets, key)
		delete(g.timers, key)
		g.mu.Unlock()
		g.flush(events)
	})
}

// Flush forces all pending buckets to be flushed immediately.
func (g *Group) Flush() {
	g.mu.Lock()
	keys := make([]string, 0, len(g.buckets))
	for k := range g.buckets {
		keys = append(keys, k)
	}
	g.mu.Unlock()

	for _, key := range keys {
		g.mu.Lock()
		t, ok := g.timers[key]
		if ok {
			t.Stop()
		}
		events := g.buckets[key]
		delete(g.buckets, key)
		delete(g.timers, key)
		g.mu.Unlock()
		g.flush(events)
	}
}

func eventKey(ev correlator.Event) string {
	return ev.Kind + ":" + ev.State.Proto + ":" + string(rune(ev.State.Port))
}
