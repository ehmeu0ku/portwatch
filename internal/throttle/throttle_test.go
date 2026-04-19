package throttle

import (
	"testing"
	"time"
)

func TestAllowFirstCallPasses(t *testing.T) {
	th := New(5 * time.Second)
	if !th.Allow("tcp:8080") {
		t.Fatal("expected first call to pass")
	}
}

func TestAllowSecondCallWithinWindowBlocked(t *testing.T) {
	th := New(5 * time.Second)
	th.Allow("tcp:8080")
	if th.Allow("tcp:8080") {
		t.Fatal("expected second call within window to be blocked")
	}
}

func TestAllowAfterWindowPasses(t *testing.T) {
	now := time.Unix(1000, 0)
	th := New(5 * time.Second)
	th.nowFn = func() time.Time { return now }

	th.Allow("tcp:8080")

	th.nowFn = func() time.Time { return now.Add(6 * time.Second) }
	if !th.Allow("tcp:8080") {
		t.Fatal("expected call after window to pass")
	}
}

func TestAllowDifferentKeysAreIndependent(t *testing.T) {
	th := New(5 * time.Second)
	th.Allow("tcp:8080")
	if !th.Allow("tcp:9090") {
		t.Fatal("expected different key to pass independently")
	}
}

func TestResetAllowsImmediateReuse(t *testing.T) {
	th := New(5 * time.Second)
	th.Allow("tcp:8080")
	th.Reset("tcp:8080")
	if !th.Allow("tcp:8080") {
		t.Fatal("expected allow after reset to pass")
	}
}

func TestResetUnknownKeyIsNoop(t *testing.T) {
	th := New(5 * time.Second)
	// Resetting a key that was never set should not panic or affect Len.
	th.Reset("tcp:8080")
	if th.Len() != 0 {
		t.Fatalf("expected 0 after resetting unknown key, got %d", th.Len())
	}
}

func TestLenTracksKeys(t *testing.T) {
	th := New(5 * time.Second)
	if th.Len() != 0 {
		t.Fatalf("expected 0, got %d", th.Len())
	}
	th.Allow("tcp:8080")
	th.Allow("tcp:9090")
	if th.Len() != 2 {
		t.Fatalf("expected 2, got %d", th.Len())
	}
}

func TestPurgeRemovesExpiredKeys(t *testing.T) {
	now := time.Unix(1000, 0)
	th := New(5 * time.Second)
	th.nowFn = func() time.Time { return now }

	th.Allow("tcp:8080")
	th.Allow("tcp:9090")

	// advance time so both entries are expired
	th.nowFn = func() time.Time { return now.Add(10 * time.Second) }
	th.Purge()

	if th.Len() != 0 {
		t.Fatalf("expected 0 after purge, got %d", th.Len())
	}
}

func TestPurgeKeepsActiveKeys(t *testing.T) {
	now := time.Unix(1000, 0)
	th := New(10 * time.Second)
	th.nowFn = func() time.Time { return now }

	th.Allow("tcp:8080") // recorded at t=1000

	th.nowFn = func() time.Time { return now.Add(6 * time.Second) }
	th.Allow("tcp:9090") // recorded at t=1006

	// advance to t=1012: 8080 is expired (12s), 9090 is still active (6s)
	th.nowFn = func() time.Time { return now.Add(12 * time.Second) }
	th.Purge()

	if th.Len() != 1 {
		t.Fatalf("expected 1 active key after purge, got %d", th.Len())
	}
}
