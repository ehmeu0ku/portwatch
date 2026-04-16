package correlator

import (
	"context"

	"github.com/user/portwatch/internal/scanner"
)

// Change pairs a Kind with a raw PortState from the monitor diff.
type Change struct {
	Kind  Kind
	State scanner.PortState
}

// Pipeline reads Changes from in, correlates each one, and sends the
// resulting PortEvent to the returned channel.  It stops when ctx is
// cancelled or in is closed.
func (c *Correlator) Pipeline(ctx context.Context, in <-chan Change) <-chan PortEvent {
	out := make(chan PortEvent, 16)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case ch, ok := <-in:
				if !ok {
					return
				}
				ev := c.Correlate(ch.Kind, ch.State)
				select {
				case out <- ev:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}
