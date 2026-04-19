package portlock

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestObserveFirstCallReturnsFalse(t *testing.T) {
	s := New(5 * time.Minute)
	if s.Observe("tcp:80", epoch) {
		t.Fatal("expected false on first observe")
	}
}

func TestObserveBeforeMinAgeReturnsFalse(t *testing.T) {
	s := New(5 * time.Minute)
	s.Observe("tcp:80", epoch)
	if s.Observe("tcp:80", epoch.Add(2*time.Minute)) {
		t.Fatal("expected false before min age")
	}
}

func TestObserveAfterMinAgeReturnsTrue(t *testing.T) {
	s := New(5 * time.Minute)
	s.Observe("tcp:80", epoch)
	if !s.Observe("tcp:80", epoch.Add(5*time.Minute)) {
		t.Fatal("expected true after min age")
	}
}

func TestIsLockedAfterPromotion(t *testing.T) {
	s := New(time.Minute)
	s.Observe("tcp:443", epoch)
	s.Observe("tcp:443", epoch.Add(time.Minute))
	if !s.IsLocked("tcp:443") {
		t.Fatal("expected port to be locked")
	}
}

func TestIsLockedUnknownKeyReturnsFalse(t *testing.T) {
	s := New(time.Minute)
	if s.IsLocked("tcp:9999") {
		t.Fatal("expected false for unknown key")
	}
}

func TestObserveAlreadyLockedReturnsFalse(t *testing.T) {
	s := New(time.Minute)
	s.Observe("tcp:22", epoch)
	s.Observe("tcp:22", epoch.Add(time.Minute)) // locks
	if s.Observe("tcp:22", epoch.Add(2*time.Minute)) {
		t.Fatal("expected false once already locked")
	}
}

func TestReleaseRemovesEntry(t *testing.T) {
	s := New(time.Minute)
	s.Observe("tcp:22", epoch)
	s.Release("tcp:22")
	if s.IsLocked("tcp:22") {
		t.Fatal("expected entry removed after release")
	}
}

func TestSnapshotReturnsCopy(t *testing.T) {
	s := New(time.Minute)
	s.Observe("tcp:80", epoch)
	snap := s.Snapshot()
	if _, ok := snap["tcp:80"]; !ok {
		t.Fatal("expected tcp:80 in snapshot")
	}
	// mutating snapshot must not affect store
	delete(snap, "tcp:80")
	if _, ok := s.Snapshot()["tcp:80"]; !ok {
		t.Fatal("store should still contain tcp:80")
	}
}
