package portbudget

import (
	"testing"
)

func TestObserveReturnsFalseUnderCeiling(t *testing.T) {
	b := New(3)
	if b.Observe("tcp") {
		t.Fatal("expected false for first observation under ceiling")
	}
	if b.Observe("tcp") {
		t.Fatal("expected false for second observation under ceiling")
	}
}

func TestObserveReturnsTrueWhenCeilingExceeded(t *testing.T) {
	b := New(2)
	b.Observe("tcp")
	b.Observe("tcp")
	if !b.Observe("tcp") {
		t.Fatal("expected true when count exceeds ceiling")
	}
}

func TestReleaseDecrementsCount(t *testing.T) {
	b := New(2)
	b.Observe("tcp")
	b.Observe("tcp")
	b.Release("tcp")
	if got := b.Count("tcp"); got != 1 {
		t.Fatalf("expected count 1 after release, got %d", got)
	}
}

func TestReleaseDoesNotGoBelowZero(t *testing.T) {
	b := New(2)
	b.Release("tcp") // should not panic or go negative
	if got := b.Count("tcp"); got != 0 {
		t.Fatalf("expected count 0, got %d", got)
	}
}

func TestExceedsReflectsCurrentState(t *testing.T) {
	b := New(1)
	b.Observe("udp")
	if b.Exceeds("udp") {
		t.Fatal("count equals ceiling, should not exceed")
	}
	b.Observe("udp")
	if !b.Exceeds("udp") {
		t.Fatal("count is over ceiling, Exceeds should return true")
	}
}

func TestDifferentProtocolsAreIndependent(t *testing.T) {
	b := New(1)
	b.Observe("tcp")
	b.Observe("tcp")
	if b.Exceeds("udp") {
		t.Fatal("udp should not be affected by tcp observations")
	}
}

func TestResetClearsCount(t *testing.T) {
	b := New(1)
	b.Observe("tcp")
	b.Observe("tcp")
	b.Reset("tcp")
	if b.Exceeds("tcp") {
		t.Fatal("expected budget to be within limit after reset")
	}
	if got := b.Count("tcp"); got != 0 {
		t.Fatalf("expected count 0 after reset, got %d", got)
	}
}

func TestZeroCeilingDisablesEnforcement(t *testing.T) {
	b := New(0)
	for i := 0; i < 100; i++ {
		if b.Observe("tcp") {
			t.Fatal("zero ceiling should never trigger Exceeds")
		}
	}
	if b.Exceeds("tcp") {
		t.Fatal("Exceeds should always return false when ceiling is zero")
	}
}

func TestSummaryWithinBudget(t *testing.T) {
	b := New(5)
	b.Observe("tcp")
	if s := b.Summary(); s != "within budget" {
		t.Fatalf("unexpected summary: %q", s)
	}
}

func TestSummaryListsBreachedProtocols(t *testing.T) {
	b := New(1)
	b.Observe("tcp")
	b.Observe("tcp")
	s := b.Summary()
	if s == "within budget" || s == "budget enforcement disabled" {
		t.Fatalf("expected breached summary, got %q", s)
	}
}
