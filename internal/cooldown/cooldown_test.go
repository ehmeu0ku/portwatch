package cooldown

import (
	"testing"
	"time"
)

func TestAllowFirstCallPasses(t *testing.T) {
	cd := New(time.Second)
	if !cd.Allow("k") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllowSecondCallWithinWindowBlocked(t *testing.T) {
	cd := New(time.Hour)
	cd.Allow("k")
	if cd.Allow("k") {
		t.Fatal("expected second call within window to be blocked")
	}
}

func TestAllowAfterWindowPasses(t *testing.T) {
	now := time.Now()
	cd := New(time.Second)
	cd.nowFunc = func() time.Time { return now }
	cd.Allow("k")
	cd.nowFunc = func() time.Time { return now.Add(2 * time.Second) }
	if !cd.Allow("k") {
		t.Fatal("expected call after window to be allowed")
	}
}

func TestAllowDifferentKeysAreIndependent(t *testing.T) {
	cd := New(time.Hour)
	cd.Allow("a")
	if !cd.Allow("b") {
		t.Fatal("expected different key to be allowed")
	}
}

func TestResetAllowsImmediateReuse(t *testing.T) {
	cd := New(time.Hour)
	cd.Allow("k")
	cd.Reset("k")
	if !cd.Allow("k") {
		t.Fatal("expected allow after reset")
	}
}

func TestLenTracksKeys(t *testing.T) {
	cd := New(time.Hour)
	if cd.Len() != 0 {
		t.Fatalf("expected 0, got %d", cd.Len())
	}
	cd.Allow("a")
	cd.Allow("b")
	if cd.Len() != 2 {
		t.Fatalf("expected 2, got %d", cd.Len())
	}
	cd.Reset("a")
	if cd.Len() != 1 {
		t.Fatalf("expected 1 after reset, got %d", cd.Len())
	}
}
