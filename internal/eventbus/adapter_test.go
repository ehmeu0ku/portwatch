package eventbus_test

import (
	"testing"

	"github.com/user/portwatch/internal/eventbus"
	"github.com/user/portwatch/internal/scanner"
)

func makePortState() scanner.PortState {
	return scanner.PortState{
		Port:    8080,
		Proto:   "tcp",
		PID:     99,
		Process: "myapp",
	}
}

func TestFromScannerStateNew(t *testing.T) {
	s := makePortState()
	e := eventbus.FromScannerState(s, eventbus.EventNew)

	if e.Type != eventbus.EventNew {
		t.Errorf("expected EventNew, got %s", e.Type)
	}
	if e.Port != s.Port || e.Proto != s.Proto || e.PID != s.PID || e.Process != s.Process {
		t.Errorf("field mismatch: event=%+v state=%+v", e, s)
	}
}

func TestFromScannerStateGone(t *testing.T) {
	s := makePortState()
	e := eventbus.FromScannerState(s, eventbus.EventGone)

	if e.Type != eventbus.EventGone {
		t.Errorf("expected EventGone, got %s", e.Type)
	}
}

func TestPublisherPublishNew(t *testing.T) {
	bus := eventbus.New()
	var got eventbus.Event
	bus.Subscribe(func(e eventbus.Event) { got = e })

	p := eventbus.NewPublisher(bus)
	p.PublishNew(makePortState())

	if got.Type != eventbus.EventNew {
		t.Errorf("expected EventNew, got %s", got.Type)
	}
	if got.Port != 8080 {
		t.Errorf("expected port 8080, got %d", got.Port)
	}
}

func TestPublisherPublishGone(t *testing.T) {
	bus := eventbus.New()
	var got eventbus.Event
	bus.Subscribe(func(e eventbus.Event) { got = e })

	p := eventbus.NewPublisher(bus)
	p.PublishGone(makePortState())

	if got.Type != eventbus.EventGone {
		t.Errorf("expected EventGone, got %s", got.Type)
	}
}
