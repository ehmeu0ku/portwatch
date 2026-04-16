package metricsink_test

import (
	"sync"
	"testing"

	"github.com/example/portwatch/internal/metricsink"
)

func TestIncAndGet(t *testing.T) {
	s := metricsink.New()
	s.Inc("hits")
	s.Inc("hits")
	if got := s.Get("hits"); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestGetMissingKeyIsZero(t *testing.T) {
	s := metricsink.New()
	if got := s.Get("nope"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestAddDelta(t *testing.T) {
	s := metricsink.New()
	s.Add("bytes", 100)
	s.Add("bytes", 50)
	if got := s.Get("bytes"); got != 150 {
		t.Fatalf("expected 150, got %d", got)
	}
}

func TestSnapshot(t *testing.T) {
	s := metricsink.New()
	s.Inc("a")
	s.Add("b", 5)
	snap := s.Snapshot()
	if snap["a"] != 1 {
		t.Fatalf("expected a=1")
	}
	if snap["b"] != 5 {
		t.Fatalf("expected b=5")
	}
}

func TestSnapshotIsCopy(t *testing.T) {
	s := metricsink.New()
	s.Inc("x")
	snap := s.Snapshot()
	s.Inc("x")
	if snap["x"] != 1 {
		t.Fatalf("snapshot should not reflect later increments")
	}
}

func TestReset(t *testing.T) {
	s := metricsink.New()
	s.Inc("c")
	s.Reset("c")
	if got := s.Get("c"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestConcurrentInc(t *testing.T) {
	s := metricsink.New()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); s.Inc("concurrent") }()
	}
	wg.Wait()
	if got := s.Get("concurrent"); got != 100 {
		t.Fatalf("expected 100, got %d", got)
	}
}
