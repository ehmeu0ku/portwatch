// Package debounce implements a key-based debouncer for portwatch.
//
// When port activity is volatile — a process repeatedly binding and releasing
// a port within a short window — raw diff events can produce noisy alerts.
// The Debouncer suppresses intermediate events by waiting for a configurable
// quiet period before invoking the downstream action.
//
// Typical usage:
//
//	d := debounce.New(500*time.Millisecond, func(key string) {
//		fmt.Println("stable event for", key)
//	})
//	d.Trigger("tcp:8080")
//	d.Trigger("tcp:8080") // resets the timer
//	// action fires ~500 ms after the last Trigger
//
// Each key is tracked independently; cancelling or re-triggering one key
// does not affect timers for other keys.
package debounce
