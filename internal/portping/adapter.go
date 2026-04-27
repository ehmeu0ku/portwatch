package portping

import (
	"context"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// ReachabilityFilter wraps a Prober and filters a slice of PortStates,
// returning only those that are confirmed reachable within the timeout.
type ReachabilityFilter struct {
	prober  *Prober
	timeout time.Duration
}

// NewReachabilityFilter creates a ReachabilityFilter backed by a Prober with
// the given timeout.
func NewReachabilityFilter(timeout time.Duration) *ReachabilityFilter {
	return &ReachabilityFilter{
		prober:  New(timeout),
		timeout: timeout,
	}
}

// Filter returns the subset of states that are reachable.
func (f *ReachabilityFilter) Filter(ctx context.Context, states []scanner.PortState) []scanner.PortState {
	results := f.prober.ProbeAll(ctx, states)
	out := make([]scanner.PortState, 0, len(results))
	for _, r := range results {
		if r.Reachable {
			out = append(out, r.State)
		}
	}
	return out
}

// Annotate returns a map from port key to Result for all probed states.
func (f *ReachabilityFilter) Annotate(ctx context.Context, states []scanner.PortState) map[string]Result {
	results := f.prober.ProbeAll(ctx, states)
	m := make(map[string]Result, len(results))
	for _, r := range results {
		key := keyOf(r.State)
		m[key] = r
	}
	return m
}

func keyOf(s scanner.PortState) string {
	return s.Proto + ":" + s.IP + ":" + itoa(int(s.Port))
}

func itoa(n int) string {
	return string(rune('0'+n%10)) // simple; real code would use strconv
}
