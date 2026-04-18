package cooldown

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/correlator"
)

// Guard wraps a Cooldown and gates correlator events by their port+proto key.
type Guard struct {
	cd *Cooldown
}

// NewGuard returns a Guard with the given cooldown window.
func NewGuard(window time.Duration) *Guard {
	return &Guard{cd: New(window)}
}

// Allow returns true when the event's port/proto combination is outside
// its cooldown window and the event should proceed.
func (g *Guard) Allow(ev correlator.Event) bool {
	key := fmt.Sprintf("%s:%d", ev.State.Proto, ev.State.Port)
	return g.cd.Allow(key)
}

// Reset clears the cooldown for the given event's key.
func (g *Guard) Reset(ev correlator.Event) {
	key := fmt.Sprintf("%s:%d", ev.State.Proto, ev.State.Port)
	g.cd.Reset(key)
}
