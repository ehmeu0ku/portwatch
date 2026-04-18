// Package cooldown provides a simple per-key cooldown tracker used to
// gate repeated actions — such as alerts or notifications — so they do
// not fire more frequently than a configured time window allows.
//
// Unlike ratelimit (which counts calls) or throttle (which uses a sliding
// window), cooldown only remembers the last allowed time and blocks until
// the full window has elapsed since that moment.
//
// Typical usage:
//
//	cd := cooldown.New(5 * time.Minute)
//	if cd.Allow(portKey) {
//		// send alert
//	}
package cooldown
