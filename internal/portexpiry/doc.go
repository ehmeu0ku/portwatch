// Package portexpiry provides a TTL-based tracker for ports that have
// disappeared from the system. Once a port has been continuously absent
// for longer than the configured TTL, Observe returns true, signalling
// that the port can be considered fully expired and removed from any
// persistent state (baselines, label stores, scorecards, etc.).
//
// Typical usage:
//
//	tracker := portexpiry.New(5 * time.Minute)
//	if tracker.Observe(key) {
//		// port has been gone long enough — clean up
//	}
//	// when port reappears:
//	tracker.Forget(key)
package portexpiry
