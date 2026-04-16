// Package circuitbreaker provides a thread-safe circuit breaker for
// protecting portwatch subsystems (e.g. notifiers, audit writers) from
// cascading failures.
//
// Usage:
//
//	br := circuitbreaker.New(3, 10*time.Second)
//
//	if err := br.Allow(); err != nil {
//		// circuit is open — skip the call
//		return err
//	}
//	if err := doWork(); err != nil {
//		br.RecordFailure()
//		return err
//	}
//	br.RecordSuccess()
//
// After `threshold` consecutive failures the breaker opens and all
// subsequent Allow calls return ErrOpen until `cooldown` has elapsed,
// at which point the next Allow call transitions it back to closed.
package circuitbreaker
