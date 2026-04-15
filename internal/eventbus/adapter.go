package eventbus

import "github.com/user/portwatch/internal/scanner"

// FromScannerState converts a scanner.PortState into an Event.
// eventType must be EventNew or EventGone.
func FromScannerState(s scanner.PortState, eventType EventType) Event {
	return Event{
		Type:    eventType,
		Port:    s.Port,
		Proto:   s.Proto,
		PID:     s.PID,
		Process: s.Process,
	}
}

// Publisher wraps a Bus and exposes helpers used by the monitor to
// emit events without importing the bus type directly.
type Publisher struct {
	bus *Bus
}

// NewPublisher returns a Publisher backed by bus.
func NewPublisher(bus *Bus) *Publisher {
	return &Publisher{bus: bus}
}

// PublishNew emits an EventNew for the given PortState.
func (p *Publisher) PublishNew(s scanner.PortState) {
	p.bus.Publish(FromScannerState(s, EventNew))
}

// PublishGone emits an EventGone for the given PortState.
func (p *Publisher) PublishGone(s scanner.PortState) {
	p.bus.Publish(FromScannerState(s, EventGone))
}
