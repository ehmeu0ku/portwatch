// Package debounce provides a simple debouncer that delays action execution
// until a quiet period has elapsed since the last trigger. This is useful for
// suppressing alert storms when ports flap rapidly.
package debounce

import (
	"sync"
	"time"
)

// Action is a function invoked after the debounce window expires.
type Action func(key string)

// Debouncer delays execution of an action until no new triggers arrive
// for a given key within the configured window.
type Debouncer struct {
	mu      sync.Mutex
	window  time.Duration
	timers  map[string]*time.Timer
	action  Action
}

// New creates a Debouncer with the given quiet window and action.
// The action is called with the key once the window expires without
// a new Trigger call for that key.
func New(window time.Duration, action Action) *Debouncer {
	return &Debouncer{
		window: window,
		timers: make(map[string]*time.Timer),
		action: action,
	}
}

// Trigger resets the debounce timer for key. If no further Trigger
// calls arrive within the window, the registered action is invoked.
func (d *Debouncer) Trigger(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Stop()
	}

	d.timers[key] = time.AfterFunc(d.window, func() {
		d.mu.Lock()
		delete(d.timers, key)
		d.mu.Unlock()
		d.action(key)
	})
}

// Cancel stops any pending timer for key without invoking the action.
func (d *Debouncer) Cancel(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Stop()
		delete(d.timers, key)
	}
}

// Pending returns the number of keys currently waiting to fire.
func (d *Debouncer) Pending() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.timers)
}
