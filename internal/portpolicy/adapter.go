package portpolicy

import (
	"github.com/example/portwatch/internal/correlator"
)

// Verdict carries the policy decision for a single event.
type Verdict struct {
	Event  correlator.Event
	Action Action
}

// Enforcer wraps a Policy and evaluates incoming correlator events.
type Enforcer struct {
	pol *Policy
}

// NewEnforcer returns an Enforcer backed by pol.
func NewEnforcer(pol *Policy) *Enforcer {
	return &Enforcer{pol: pol}
}

// Evaluate returns a Verdict for the given event.
func (e *Enforcer) Evaluate(ev correlator.Event) Verdict {
	action := e.pol.Evaluate(ev.State.Port, ev.State.Proto)
	return Verdict{Event: ev, Action: action}
}

// IsDenied is a convenience helper.
func (e *Enforcer) IsDenied(ev correlator.Event) bool {
	return e.Evaluate(ev).Action == Deny
}
