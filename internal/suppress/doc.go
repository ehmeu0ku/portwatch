// Package suppress implements event suppression for portwatch.
//
// When a port appears or disappears, the monitor may fire the same
// alert repeatedly across scan cycles. The Suppressor ensures that
// identical events are forwarded at most once per configured window,
// reducing alert fatigue without losing the initial notification.
//
// Usage:
//
//	s := suppress.New(5 * time.Minute)
//
//	key := suppress.Key{Proto: "tcp", Addr: "0.0.0.0", Port: 8080, Kind: "new"}
//	if s.Allow(key) {
//		// forward the alert
//	}
//
// Call Purge() periodically (e.g. from a background goroutine) to
// reclaim memory for keys whose suppression window has elapsed.
package suppress
