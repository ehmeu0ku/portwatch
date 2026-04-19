package portexpiry

import (
	"testing"
	"time"
)

func TestObserveFirstCallReturnsFalse(t *testing.T) {
	tr := New(5 * time.Minute)
	if tr.Observe("tcp:8080") {
		t.Fatal("expected false on first observation")
	}
}

func TestObserveWithinTTLReturnsFalse(t *testing.T) {
	now := time.Now()
	tr := New(5 * time.Minute)
	tr.now = func() time.Time { return now }
	tr.Observe("tcp:8080")
	tr.now = func() time.Time { return now.Add(4 * time.Minute) }
	if tr.Observe("tcp:8080") {
		t.Fatal("expected false before TTL expires")
	}
}

func TestObserveAfterTTLReturnsTrue(t *testing.T) {
	now := time.Now()
	tr := New(5 * time.Minute)
	tr.now = func() time.Time { return now }
	tr.Observe("tcp:8080")
	tr.now = func() time.Time { return now.Add(5 * time.Minute) }
	if !tr.Observe("tcp:8080") {
		t.Fatal("expected true after TTL expires")
	}
}

func TestForgetResetsTracking(t *testing.T) {
	now := time.Now()
	tr := New(1 * time.Minute)
	tr.now = func() time.Time { return now }
	tr.Observe("tcp:9090")
	tr.now = func() time.Time { return now.Add(2 * time.Minute) }
	tr.Forget("tcp:9090")
	// after forget, first observation resets the clock
	tr.now = func() time.Time { return now.Add(2 * time.Minute) }
	if tr.Observe("tcp:9090") {
		t.Fatal("expected false after forget resets entry")
	}
}

func TestLenReflectsTrackedCount(t *testing.T) {
	tr := New(time.Minute)
	tr.Observe("tcp:80")
	tr.Observe("tcp:443")
	if tr.Len() != 2 {
		t.Fatalf("expected 2, got %d", tr.Len())
	}
	tr.Forget("tcp:80")
	if tr.Len() != 1 {
		t.Fatalf("expected 1, got %d", tr.Len())
	}
}

func TestDifferentKeysAreIndependent(t *testing.T) {
	now := time.Now()
	tr := New(time.Minute)
	tr.now = func() time.Time { return now }
	tr.Observe("tcp:80")
	tr.now = func() time.Time { return now.Add(2 * time.Minute) }
	tr.Observe("tcp:443") // first observation for this key
	if !tr.Observe("tcp:80") {
		t.Fatal("tcp:80 should have expired")
	}
	if tr.Observe("tcp:443") {
		t.Fatal("tcp:443 should not have expired yet")
	}
}
