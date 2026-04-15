// Package throttle provides per-key event throttling used by portwatch to
// suppress repeated alerts for the same port within a configurable time window.
//
// # Overview
//
// When a new or unexpected port is detected, the monitor may fire on every
// scan tick until the port disappears. Throttle ensures that only the first
// detection within each window is forwarded to the alerter, reducing noise
// for long-lived unexpected listeners.
//
// # Usage
//
//	th := throttle.New(30 * time.Second)
//
//	if th.Allow("tcp:8080") {
//	    alerter.Notify(state)
//	}
//
// Call Purge() periodically (e.g. from the monitor tick) to reclaim memory
// for ports that are no longer being tracked.
package throttle
