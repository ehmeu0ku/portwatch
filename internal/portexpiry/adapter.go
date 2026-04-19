package portexpiry

import (
	"fmt"

	"github.com/user/portwatch/internal/scanner"
)

// KeyFromState returns a canonical string key for a PortState.
func KeyFromState(s scanner.PortState) string {
	return fmt.Sprintf("%s:%d", s.Proto, s.Port)
}

// Observer wraps Tracker and provides scanner-aware helpers.
type Observer struct {
	tracker *Tracker
}

// NewObserver returns an Observer backed by the given Tracker.
func NewObserver(tr *Tracker) *Observer {
	return &Observer{tracker: tr}
}

// ObserveGone records that the port described by s is absent.
// Returns true when the port has been absent longer than the TTL.
func (o *Observer) ObserveGone(s scanner.PortState) bool {
	return o.tracker.Observe(KeyFromState(s))
}

// ObservePresent resets any absence record for s (port came back).
func (o *Observer) ObservePresent(s scanner.PortState) {
	o.tracker.Forget(KeyFromState(s))
}
