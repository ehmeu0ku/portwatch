package portdrift

import (
	"github.com/user/portwatch/internal/eventbus"
	"github.com/user/portwatch/internal/scanner"
)

// DriftKind is the event kind published on the bus when address drift is detected.
const DriftKind = "port.drift"

// Publisher wraps a Detector and publishes DriftEvents onto an eventbus.
type Publisher struct {
	detector *Detector
	bus      *eventbus.Bus
}

// NewPublisher returns a Publisher that uses det to detect drift and
// publishes any events onto bus.
func NewPublisher(det *Detector, bus *eventbus.Bus) *Publisher {
	return &Publisher{detector: det, bus: bus}
}

// Observe runs drift detection over states and publishes one bus event per
// DriftEvent detected. The payload of each published event is the DriftEvent
// value itself.
func (p *Publisher) Observe(states []scanner.PortState) {
	events := p.detector.Observe(states)
	for _, e := range events {
		e := e // capture
		p.bus.Publish(eventbus.Event{
			Kind:    DriftKind,
			Payload: e,
		})
	}
}
