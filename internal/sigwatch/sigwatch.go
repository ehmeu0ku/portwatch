// Package sigwatch listens for OS signals and publishes shutdown or reload
// events onto the event bus so other components can react cleanly.
package sigwatch

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/eventbus"
)

// Handler subscribes to OS signals and forwards them as lifecycle events.
type Handler struct {
	bus    *eventbus.Bus
	signals []os.Signal
}

// New creates a Handler that watches the given signals.
// If no signals are provided it defaults to SIGINT, SIGTERM, and SIGHUP.
func New(bus *eventbus.Bus, sigs ...os.Signal) *Handler {
	if len(sigs) == 0 {
		sigs = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP}
	}
	return &Handler{bus: bus, signals: sigs}
}

// Run blocks until ctx is cancelled or a watched signal arrives.
// SIGHUP publishes a "reload" event; all others publish "shutdown".
func (h *Handler) Run(ctx context.Context) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, h.signals...)
	defer signal.Stop(ch)

	for {
		select {
		case <-ctx.Done():
			return
		case sig := <-ch:
			kind := "shutdown"
			if sig == syscall.SIGHUP {
				kind = "reload"
			}
			h.bus.Publish(eventbus.Event{
				Kind:    kind,
				Payload: sig.String(),
			})
			if kind == "shutdown" {
				return
			}
		}
	}
}
