package portlifecycle

import (
	"testing"
	"time"
)

func TestObserveFirstCallSetsFirstSeen(t *testing.T) {
	tr := New()
	before := time.Now()
	e := tr.Observe("tcp:8080")
	if e.FirstSeen.Before(before) {
		t.Fatalf("expected FirstSeen >= %v, got %v", before, e.FirstSeen)
	}
	if e.SeenCount != 1 {
		t.Fatalf("expected SeenCount=1, got %d", e.SeenCount)
	}
}

func TestObserveIncrementsSeenCount(t *testing.T) {
	tr := New()
	tr.Observe("tcp:9090")
	tr.Observe("tcp:9090")
	e := tr.Observe("tcp:9090")
	if e.SeenCount != 3 {
		t.Fatalf("expected SeenCount=3, got %d", e.SeenCount)
	}
}

func TestObservePreservesFirstSeen(t *testing.T) {
	tr := New()
	first := tr.Observe("tcp:443")
	time.Sleep(2 * time.Millisecond)
	second := tr.Observe("tcp:443")
	if !second.FirstSeen.Equal(first.FirstSeen) {
		t.Fatalf("FirstSeen changed: %v -> %v", first.FirstSeen, second.FirstSeen)
	}
}

func TestMarkGoneSetsGoneAt(t *testing.T) {
	tr := New()
	tr.Observe("udp:53")
	before := time.Now()
	tr.MarkGone("udp:53")
	e, ok := tr.Get("udp:53")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.GoneAt == nil {
		t.Fatal("expected GoneAt to be set")
	}
	if e.GoneAt.Before(before) {
		t.Fatalf("GoneAt too early: %v", e.GoneAt)
	}
}

func TestObserveAfterGoneClearsGoneAt(t *testing.T) {
	tr := New()
	tr.Observe("tcp:22")
	tr.MarkGone("tcp:22")
	e := tr.Observe("tcp:22")
	if e.GoneAt != nil {
		t.Fatal("expected GoneAt to be cleared after re-observe")
	}
}

func TestMarkGoneUnknownKeyIsNoop(t *testing.T) {
	tr := New()
	tr.MarkGone("tcp:9999") // should not panic
	if tr.Len() != 0 {
		t.Fatalf("expected empty tracker, got len=%d", tr.Len())
	}
}

func TestGetMissingKeyReturnsFalse(t *testing.T) {
	tr := New()
	_, ok := tr.Get("tcp:1234")
	if ok {
		t.Fatal("expected ok=false for missing key")
	}
}

func TestForgetRemovesEntry(t *testing.T) {
	tr := New()
	tr.Observe("tcp:80")
	tr.Forget("tcp:80")
	_, ok := tr.Get("tcp:80")
	if ok {
		t.Fatal("expected entry to be removed after Forget")
	}
	if tr.Len() != 0 {
		t.Fatalf("expected Len=0, got %d", tr.Len())
	}
}

func TestLenReflectsTrackedCount(t *testing.T) {
	tr := New()
	tr.Observe("tcp:80")
	tr.Observe("tcp:443")
	tr.Observe("udp:53")
	if tr.Len() != 3 {
		t.Fatalf("expected Len=3, got %d", tr.Len())
	}
}

func TestEntryStringWithoutGone(t *testing.T) {
	e := &Entry{
		FirstSeen: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		LastSeen:  time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		SeenCount: 5,
	}
	s := e.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}
