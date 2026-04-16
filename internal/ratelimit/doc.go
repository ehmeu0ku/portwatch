// Package ratelimit provides a simple per-key rate limiter that blocks
// repeated events within a configurable cooldown window.
//
// Usage:
//
//	rl := ratelimit.New(5 * time.Second)
//	if rl.Allow("port:8080") {
//		// handle event
//	}
//
// Each key is tracked independently. Calling Reset on a key allows it to
// fire immediately on the next Allow call.
package ratelimit
