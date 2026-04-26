package portmap

import (
	"context"

	"github.com/user/portwatch/internal/scanner"
)

// SyncFunc is called whenever the port map is updated with a reconciled
// snapshot of all currently-active port states.
type SyncFunc func(states []scanner.PortState)

// Syncer keeps a PortMap in sync with a stream of scanner states produced
// by the monitor diff loop.  It is the glue between the raw diff events
// (new / gone) and the authoritative in-memory map of what is currently
// listening.
type Syncer struct {
	pm       *PortMap
	onUpdate SyncFunc
}

// NewSyncer returns a Syncer that maintains pm and calls onUpdate (if
// non-nil) after every successful reconciliation.
func NewSyncer(pm *PortMap, onUpdate SyncFunc) *Syncer {
	return &Syncer{pm: pm, onUpdate: onUpdate}
}

// HandleNew records a newly-observed port state into the map.  It is safe
// to call concurrently.
func (s *Syncer) HandleNew(state scanner.PortState) {
	s.pm.Set(KeyFromState(state), state)
	s.notify()
}

// HandleGone removes a port state that is no longer observed.  It is safe
// to call concurrently.
func (s *Syncer) HandleGone(state scanner.PortState) {
	s.pm.Delete(KeyFromState(state))
	s.notify()
}

// Snapshot returns a point-in-time copy of all states currently held in
// the underlying PortMap.
func (s *Syncer) Snapshot() []scanner.PortState {
	return s.pm.All()
}

// Run reads new/gone events from the provided channels until ctx is
// cancelled or both channels are closed.  It is intended for use in a
// dedicated goroutine.
//
//	newCh  – receives states that just appeared
//	goneCh – receives states that just disappeared
func (s *Syncer) Run(ctx context.Context, newCh, goneCh <-chan scanner.PortState) {
	for {
		select {
		case <-ctx.Done():
			return
		case st, ok := <-newCh:
			if !ok {
				return
			}
			s.HandleNew(st)
		case st, ok := <-goneCh:
			if !ok {
				return
			}
			s.HandleGone(st)
		}
	}
}

// notify calls the registered SyncFunc with the current snapshot, if one
// was provided at construction time.
func (s *Syncer) notify() {
	if s.onUpdate == nil {
		return
	}
	s.onUpdate(s.pm.All())
}
