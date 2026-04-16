package rollup_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/correlator"
	"github.com/user/portwatch/internal/rollup"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(kind string, port uint16) correlator.Event {
	return correlator.Event{
		Kind:  kind,
		State: scanner.PortState{Port: port, Proto: "tcp"},
	}
}

func TestSingleEventFlushedAfterWindow(t *testing.T) {
	var mu sync.Mutex
	var got [][]correlator.Event

	g := rollup.New(50*time.Millisecond, func(evs []correlator.Event) {
		mu.Lock()
		got = append(got, evs)
		mu.Unlock()
	})

	g.Add(makeEvent("new", 8080))
	time.Sleep(120 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(got) != 1 || len(got[0]) != 1 {
		t.Fatalf("expected 1 batch of 1 event, got %v", got)
	}
}

func TestBurstGroupedIntoOneBatch(t *testing.T) {
	var mu sync.Mutex
	var got [][]correlator.Event

	g := rollup.New(60*time.Millisecond, func(evs []correlator.Event) {
		mu.Lock()
		got = append(got, evs)
		mu.Unlock()
	})

	for i := 0; i < 5; i++ {
		g.Add(makeEvent("new", 9000))
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(120 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(got) != 1 {
		t.Fatalf("expected 1 batch, got %d", len(got))
	}
	if len(got[0]) != 5 {
		t.Fatalf("expected 5 events in batch, got %d", len(got[0]))
	}
}

func TestDifferentKeysFlushIndependently(t *testing.T) {
	var mu sync.Mutex
	var got [][]correlator.Event

	g := rollup.New(50*time.Millisecond, func(evs []correlator.Event) {
		mu.Lock()
		got = append(got, evs)
		mu.Unlock()
	})

	g.Add(makeEvent("new", 8080))
	g.Add(makeEvent("new", 9090))
	time.Sleep(120 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(got) != 2 {
		t.Fatalf("expected 2 batches, got %d", len(got))
	}
}

func TestFlushDrainsAllBuckets(t *testing.T) {
	var mu sync.Mutex
	var got [][]correlator.Event

	g := rollup.New(10*time.Second, func(evs []correlator.Event) {
		mu.Lock()
		got = append(got, evs)
		mu.Unlock()
	})

	g.Add(makeEvent("new", 8080))
	g.Add(makeEvent("gone", 443))
	g.Flush()

	mu.Lock()
	defer mu.Unlock()
	if len(got) != 2 {
		t.Fatalf("expected 2 batches after Flush, got %d", len(got))
	}
}
