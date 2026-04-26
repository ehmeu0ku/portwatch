package portjitter

import (
	"testing"
	"time"
)

func TestObserveFirstCallReturnsFalse(t *testing.T) {
	tr := New(5)
	now := time.Now()
	_, stable := tr.Observe("tcp:8080", now)
	if stable {
		t.Fatal("expected stable=false on first observation")
	}
}

func TestObserveSecondCallReturnsFalse(t *testing.T) {
	tr := New(5)
	now := time.Now()
	tr.Observe("tcp:8080", now)
	_, stable := tr.Observe("tcp:8080", now.Add(time.Second))
	if stable {
		t.Fatal("expected stable=false with only one interval")
	}
}

func TestObserveThirdCallReturnsStable(t *testing.T) {
	tr := New(5)
	now := time.Now()
	tr.Observe("tcp:8080", now)
	tr.Observe("tcp:8080", now.Add(time.Second))
	j, stable := tr.Observe("tcp:8080", now.Add(2*time.Second))
	if !stable {
		t.Fatal("expected stable=true with two intervals")
	}
	// equal intervals → zero spread
	if j != 0 {
		t.Fatalf("expected zero jitter for equal intervals, got %v", j)
	}
}

func TestObserveHighJitterDetected(t *testing.T) {
	tr := New(10)
	now := time.Now()
	tr.Observe("tcp:9000", now)
	tr.Observe("tcp:9000", now.Add(100*time.Millisecond))
	j, stable := tr.Observe("tcp:9000", now.Add(900*time.Millisecond))
	if !stable {
		t.Fatal("expected stable=true")
	}
	if j == 0 {
		t.Fatal("expected non-zero jitter for unequal intervals")
	}
}

func TestForgetRemovesKey(t *testing.T) {
	tr := New(5)
	now := time.Now()
	tr.Observe("tcp:443", now)
	tr.Forget("tcp:443")
	if tr.Len() != 0 {
		t.Fatalf("expected 0 entries after forget, got %d", tr.Len())
	}
}

func TestForgetThenObserveResetsHistory(t *testing.T) {
	tr := New(5)
	now := time.Now()
	tr.Observe("tcp:443", now)
	tr.Observe("tcp:443", now.Add(time.Second))
	tr.Forget("tcp:443")
	tr.Observe("tcp:443", now.Add(2*time.Second))
	_, stable := tr.Observe("tcp:443", now.Add(3*time.Second))
	// only one interval after reset
	if stable {
		t.Fatal("expected stable=false after forget and two observations")
	}
}

func TestWindowEvictsOldIntervals(t *testing.T) {
	tr := New(2)
	now := time.Now()
	tr.Observe("tcp:22", now)
	tr.Observe("tcp:22", now.Add(1*time.Second))
	tr.Observe("tcp:22", now.Add(2*time.Second))
	tr.Observe("tcp:22", now.Add(3*time.Second))
	_, stable := tr.Observe("tcp:22", now.Add(4*time.Second))
	if !stable {
		t.Fatal("expected stable=true")
	}
	if tr.entries["tcp:22"] == nil {
		t.Fatal("entry should exist")
	}
	if len(tr.entries["tcp:22"].intervals) > 2 {
		t.Fatalf("expected at most 2 intervals, got %d", len(tr.entries["tcp:22"].intervals))
	}
}

func TestDefaultWindowAppliedWhenZero(t *testing.T) {
	tr := New(0)
	if tr.window != 10 {
		t.Fatalf("expected default window 10, got %d", tr.window)
	}
}
