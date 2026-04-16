package correlator_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/correlator"
	"github.com/user/portwatch/internal/scanner"
)

func TestPipelineDeliversEvent(t *testing.T) {
	c := newCorrelator()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	in := make(chan correlator.Change, 1)
	out := c.Pipeline(ctx, in)

	in <- correlator.Change{Kind: correlator.KindNew, State: scanner.PortState{Port: 8080, Proto: "tcp"}}
	close(in)

	select {
	case ev, ok := <-out:
		if !ok {
			t.Fatal("channel closed before event received")
		}
		if ev.Kind != correlator.KindNew {
			t.Fatalf("expected KindNew, got %s", ev.Kind)
		}
		if ev.State.Port != 8080 {
			t.Fatalf("expected port 8080, got %d", ev.State.Port)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for event")
	}
}

func TestPipelineClosesOnContextCancel(t *testing.T) {
	c := newCorrelator()
	ctx, cancel := context.WithCancel(context.Background())

	in := make(chan correlator.Change)
	out := c.Pipeline(ctx, in)

	cancel()

	select {
	case _, ok := <-out:
		if ok {
			t.Fatal("expected channel to be closed after cancel")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for pipeline to stop")
	}
}

func TestPipelineMultipleChanges(t *testing.T) {
	c := newCorrelator()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	in := make(chan correlator.Change, 3)
	out := c.Pipeline(ctx, in)

	ports := []uint16{22, 80, 443}
	for _, p := range ports {
		in <- correlator.Change{Kind: correlator.KindNew, State: scanner.PortState{Port: p, Proto: "tcp"}}
	}
	close(in)

	received := 0
	for range out {
		received++
	}
	if received != len(ports) {
		t.Fatalf("expected %d events, got %d", len(ports), received)
	}
}
