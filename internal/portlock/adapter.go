package portlock

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/eventbus"
	"github.com/user/portwatch/internal/scanner"
)

// KeyFromState returns the canonical store key for a PortState.
func KeyFromState(ps scanner.PortState) string {
	return fmt.Sprintf("%s:%d", ps.Proto, ps.Port)
}

// Observer wraps Store and integrates with the event bus, observing
// NEW events and releasing on GONE events.
type Observer struct {
	store *Store
	now   func() time.Time
}

// NewObserver creates an Observer backed by store.
func NewObserver(store *Store) *Observer {
	return &Observer{store: store, now: time.Now}
}

// Handle processes an eventbus.Event, updating lock state accordingly.
// It returns true if the event caused a port to become locked.
func (o *Observer) Handle(ev eventbus.Event) bool {
	key := fmt.Sprintf("%s:%d", ev.State.Proto, ev.State.Port)
	switch ev.Kind {
	case eventbus.KindNew:
		return o.store.Observe(key, o.now())
	case eventbus.KindGone:
		o.store.Release(key)
	}
	return false
}

// IsLocked reports whether the port described by ps is locked.
func (o *Observer) IsLocked(ps scanner.PortState) bool {
	return o.store.IsLocked(KeyFromState(ps))
}
