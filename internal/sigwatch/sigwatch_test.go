package sigwatch_test

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/user/portwatch/internal/eventbus"
	"github.com/user/portwatch/internal/sigwatch"
)

func TestShutdownSignalPublishesEvent(t *testing.T) {
	bus := eventbus.New()
	var got eventbus.Event
	bus.Subscribe("test", func(e eventbus.Event) { got = e })

	h := sigwatch.New(bus, syscall.SIGUSR1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() { h.Run(ctx); close(done) }()

	time.Sleep(20 * time.Millisecond)
	syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("handler did not return after signal")
	}

	if got.Kind != "shutdown" {
		t.Fatalf("expected shutdown, got %q", got.Kind)
	}
}

func TestReloadSignalPublishesReloadEvent(t *testing.T) {
	bus := eventbus.New()
	events := make([]eventbus.Event, 0)
	bus.Subscribe("test", func(e eventbus.Event) { events = append(events, e) })

	h := sigwatch.New(bus, syscall.SIGHUP)
	ctx, cancel := context.WithCancel(context.Background())

	go h.Run(ctx)

	time.Sleep(20 * time.Millisecond)
	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	time.Sleep(50 * time.Millisecond)
	cancel()
	time.Sleep(20 * time.Millisecond)

	if len(events) == 0 {
		t.Fatal("expected at least one event")
	}
	if events[0].Kind != "reload" {
		t.Fatalf("expected reload, got %q", events[0].Kind)
	}
}

func TestContextCancelStopsHandler(t *testing.T) {
	bus := eventbus.New()
	h := sigwatch.New(bus)
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() { h.Run(ctx); close(done) }()

	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("handler did not stop after context cancel")
	}
}
