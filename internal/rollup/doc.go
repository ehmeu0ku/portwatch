// Package rollup provides burst-grouping for port events.
//
// During service restarts or network reconfiguration, portwatch may observe
// dozens of open/close events for the same port in quick succession.  Sending
// one alert per event would overwhelm operators.
//
// rollup.Group collects events that share the same (kind, proto, port) key
// within a configurable sliding window.  When the window expires without a
// new event the whole batch is forwarded to a caller-supplied flush function,
// which can format and dispatch a single consolidated notification.
//
// Usage:
//
//	g := rollup.New(2*time.Second, func(evs []correlator.Event) {
//		fmt.Printf("burst of %d events on port %d\n", len(evs), evs[0].State.Port)
//	})
//	g.Add(ev)
package rollup
