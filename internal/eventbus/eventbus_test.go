package eventbus_test

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/eventbus"
)

func makeEvent(t EventType) eventbus.Event {
	return eventbus.Event{Type: t, Port: 9090, Proto: "tcp", PID: 1, Process: "svc"}
}

type EventType = eventbus.EventType

func TestSubscribeReceivesPublishedEvent(t *testing.T) {
	bus := eventbus.New()
	var got []eventbus.Event

	bus.Subscribe(func(e eventbus.Event) {
		got = append(got, e)
	})

	bus.Publish(makeEvent(eventbus.EventNew))

	if len(got) != 1 {
		t.Fatalf("expected 1 event, got %d", len(got))
	}
	if got[0].Type != eventbus.EventNew {
		t.Errorf("expected EventNew, got %s", got[0].Type)
	}
}

func TestMultipleSubscribersAllReceive(t *testing.T) {
	bus := eventbus.New()
	var mu sync.Mutex
	count := 0

	for i := 0; i < 3; i++ {
		bus.Subscribe(func(e eventbus.Event) {
			mu.Lock()
			count++
			mu.Unlock()
		})
	}

	bus.Publish(makeEvent(eventbus.EventGone))

	if count != 3 {
		t.Errorf("expected 3 deliveries, got %d", count)
	}
}

func TestUnsubscribeStopsDelivery(t *testing.T) {
	bus := eventbus.New()
	var got []eventbus.Event

	unsub := bus.Subscribe(func(e eventbus.Event) {
		got = append(got, e)
	})

	unsub()
	bus.Publish(makeEvent(eventbus.EventNew))

	if len(got) != 0 {
		t.Errorf("expected no events after unsubscribe, got %d", len(got))
	}
}

func TestLenReflectsActiveHandlers(t *testing.T) {
	bus := eventbus.New()

	if bus.Len() != 0 {
		t.Fatalf("expected 0 handlers initially")
	}

	unsub1 := bus.Subscribe(func(e eventbus.Event) {})
	bus.Subscribe(func(e eventbus.Event) {})

	if bus.Len() != 2 {
		t.Errorf("expected 2 handlers, got %d", bus.Len())
	}

	unsub1()

	if bus.Len() != 1 {
		t.Errorf("expected 1 handler after unsubscribe, got %d", bus.Len())
	}
}

func TestPublishWithNoSubscribersIsSafe(t *testing.T) {
	bus := eventbus.New()
	// Should not panic.
	bus.Publish(makeEvent(eventbus.EventNew))
}

func TestEventFieldsPreserved(t *testing.T) {
	bus := eventbus.New()
	var received eventbus.Event

	bus.Subscribe(func(e eventbus.Event) { received = e })

	sent := eventbus.Event{Type: eventbus.EventGone, Port: 443, Proto: "tcp", PID: 42, Process: "nginx"}
	bus.Publish(sent)

	if received != sent {
		t.Errorf("received event %+v does not match sent %+v", received, sent)
	}
}
