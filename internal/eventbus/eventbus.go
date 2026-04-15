// Package eventbus provides a simple publish/subscribe mechanism for
// broadcasting port change events to multiple consumers within portwatch.
package eventbus

import (
	"sync"
)

// EventType classifies a port change notification.
type EventType string

const (
	EventNew  EventType = "NEW"
	EventGone EventType = "GONE"
)

// Event carries a single port-change notification.
type Event struct {
	Type    EventType
	Port    uint16
	Proto   string
	PID     int
	Process string
}

// Handler is a callback invoked for each published event.
type Handler func(Event)

// Bus is a thread-safe publish/subscribe event bus.
type Bus struct {
	mu       sync.RWMutex
	handlers []Handler
}

// New returns an initialised, empty Bus.
func New() *Bus {
	return &Bus{}
}

// Subscribe registers h to receive all future events.
// Returns an unsubscribe function that removes the handler.
func (b *Bus) Subscribe(h Handler) func() {
	b.mu.Lock()
	defer b.mu.Unlock()

	idx := len(b.handlers)
	b.handlers = append(b.handlers, h)

	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		b.handlers[idx] = nil
	}
}

// Publish delivers e to all registered (non-nil) handlers.
// Each handler is called synchronously in the order it was registered.
func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, h := range b.handlers {
		if h != nil {
			h(e)
		}
	}
}

// Len returns the number of registered (non-nil) handlers.
func (b *Bus) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	count := 0
	for _, h := range b.handlers {
		if h != nil {
			count++
		}
	}
	return count
}
