package graceperiod

import (
	"testing"
	"time"
)

func TestObserveFirstCallReturnsFalse(t *testing.T) {
	f := New(100 * time.Millisecond)
	if f.Observe("tcp:8080") {
		t.Fatal("expected false on first observe")
	}
}

func TestObserveWithinWindowReturnsFalse(t *testing.T) {
	f := New(1 * time.Hour)
	f.Observe("tcp:8080")
	if f.Observe("tcp:8080") {
		t.Fatal("expected false within grace window")
	}
}

func TestObserveAfterWindowReturnsTrue(t *testing.T) {
	now := time.Now()
	f := New(50 * time.Millisecond)
	f.now = func() time.Time { return now }
	f.Observe("tcp:9090")

	// advance clock beyond window
	f.now = func() time.Time { return now.Add(100 * time.Millisecond) }
	if !f.Observe("tcp:9090") {
		t.Fatal("expected true after grace window elapsed")
	}
}

func TestForgetRemovesKey(t *testing.T) {
	f := New(1 * time.Hour)
	f.Observe("tcp:443")
	if f.Len() != 1 {
		t.Fatalf("expected 1 pending, got %d", f.Len())
	}
	f.Forget("tcp:443")
	if f.Len() != 0 {
		t.Fatalf("expected 0 pending after forget, got %d", f.Len())
	}
}

func TestForgetThenObserveResetsTimer(t *testing.T) {
	now := time.Now()
	f := New(50 * time.Millisecond)
	f.now = func() time.Time { return now }
	f.Observe("tcp:22")

	f.now = func() time.Time { return now.Add(100 * time.Millisecond) }
	f.Forget("tcp:22")

	// re-observe after forget should reset the clock
	if f.Observe("tcp:22") {
		t.Fatal("expected false after forget resets timer")
	}
}

func TestLenTracksMultipleKeys(t *testing.T) {
	f := New(1 * time.Hour)
	f.Observe("tcp:80")
	f.Observe("tcp:443")
	f.Observe("udp:53")
	if f.Len() != 3 {
		t.Fatalf("expected 3, got %d", f.Len())
	}
}
