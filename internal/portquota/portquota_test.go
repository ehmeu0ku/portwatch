package portquota

import (
	"testing"
)

func TestObserveReturnsFalseUnderCeiling(t *testing.T) {
	q := New(3)
	if q.Observe("tcp") {
		t.Fatal("expected false when count <= ceiling")
	}
}

func TestObserveReturnsTrueWhenCeilingExceeded(t *testing.T) {
	q := New(2)
	q.Observe("tcp")
	q.Observe("tcp")
	if !q.Observe("tcp") {
		t.Fatal("expected true when count exceeds ceiling")
	}
}

func TestExceedsReflectsCurrentState(t *testing.T) {
	q := New(1)
	q.Observe("udp")
	if q.Exceeds("udp") {
		t.Fatal("should not exceed ceiling at count==ceiling")
	}
	q.Observe("udp")
	if !q.Exceeds("udp") {
		t.Fatal("should exceed ceiling when count > ceiling")
	}
}

func TestReleaseDecrementsCount(t *testing.T) {
	q := New(1)
	q.Observe("tcp")
	q.Observe("tcp") // now exceeds
	q.Release("tcp")
	if q.Exceeds("tcp") {
		t.Fatal("expected count to drop below ceiling after release")
	}
}

func TestReleaseDoesNotGoBelowZero(t *testing.T) {
	q := New(5)
	q.Release("tcp") // no-op
	snap := q.Snapshot()
	if e, ok := snap["tcp"]; ok && e.Count < 0 {
		t.Fatalf("count went negative: %d", e.Count)
	}
}

func TestSetCeilingOverridesDefault(t *testing.T) {
	q := New(10)
	q.SetCeiling("tcp", 2)
	q.Observe("tcp")
	q.Observe("tcp")
	if !q.Observe("tcp") {
		t.Fatal("expected ceiling breach with explicit ceiling of 2")
	}
}

func TestDifferentProtocolsAreIndependent(t *testing.T) {
	q := New(1)
	q.Observe("tcp")
	q.Observe("tcp") // tcp exceeds
	if q.Exceeds("udp") {
		t.Fatal("udp should not be affected by tcp count")
	}
}

func TestSnapshotReturnsCopy(t *testing.T) {
	q := New(5)
	q.Observe("tcp")
	snap := q.Snapshot()
	snap["tcp"] = Entry{Count: 99, Ceiling: 99} // mutate copy
	if q.Snapshot()["tcp"].Count == 99 {
		t.Fatal("snapshot mutation affected internal state")
	}
}

func TestSnapshotCeilingMatchesExplicit(t *testing.T) {
	q := New(10)
	q.SetCeiling("udp", 7)
	q.Observe("udp")
	snap := q.Snapshot()
	if snap["udp"].Ceiling != 7 {
		t.Fatalf("expected ceiling 7, got %d", snap["udp"].Ceiling)
	}
}
