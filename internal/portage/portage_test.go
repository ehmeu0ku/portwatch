package portage

import (
	"testing"
	"time"
)

func TestObserveFirstCallSetsFirstSeen(t *testing.T) {
	tr := New()
	fixed := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	tr.now = func() time.Time { return fixed }

	e := tr.Observe("tcp:8080")
	if !e.FirstSeen.Equal(fixed) {
		t.Fatalf("expected FirstSeen %v, got %v", fixed, e.FirstSeen)
	}
}

func TestObserveSecondCallPreservesFirstSeen(t *testing.T) {
	tr := New()
	t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := t1.Add(5 * time.Minute)
	call := 0
	tr.now = func() time.Time {
		call++
		if call == 1 {
			return t1
		}
		return t2
	}

	tr.Observe("tcp:8080")
	e := tr.Observe("tcp:8080")

	if !e.FirstSeen.Equal(t1) {
		t.Fatalf("expected FirstSeen %v, got %v", t1, e.FirstSeen)
	}
	if !e.LastSeen.Equal(t2) {
		t.Fatalf("expected LastSeen %v, got %v", t2, e.LastSeen)
	}
}

func TestAgeCalculation(t *testing.T) {
	tr := New()
	t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	tr.now = func() time.Time { return t1 }
	443")

	e, _ := tr.Get("tcp:443")
	now := t1.Add(10 * time.Minute)
	if e.Age(now) != 10*time.Minute {
		t.Fatalf("expected 10m age, got %v", e.Age(now))
	}
}

func TestForgetRemovesEntry(t *testing.T) {
	tr := New()
	tr.Observe("tcp:22")
	tr.Forget("tcp:22")
	_, ok := tr.Get("tcp:22")
	if ok {
		t.Fatal("expected entry to be removed")
	}
}

func TestGetMissingReturnsFalse(t *testing.T) {
	tr := New()
	_, ok := tr.Get("tcp:9999")
	if ok {
		t.Fatal("expected false for missing key")
	}
}

func TestLenReflectsObservations(t *testing.T) {
	tr := New()
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
