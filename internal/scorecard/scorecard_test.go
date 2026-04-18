package scorecard

import (
	"testing"
)

func TestRecordIncrementsSeenCount(t *testing.T) {
	sc := New()
	sc.Record("tcp:8080", false)
	sc.Record("tcp:8080", false)
	e, ok := sc.Get("tcp:8080")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.SeenCount != 2 {
		t.Fatalf("expected SeenCount 2, got %d", e.SeenCount)
	}
}

func TestRecordIncrementsAlertCount(t *testing.T) {
	sc := New()
	sc.Record("tcp:8080", true)
	sc.Record("tcp:8080", false)
	e, _ := sc.Get("tcp:8080")
	if e.AlertCount != 1 {
		t.Fatalf("expected AlertCount 1, got %d", e.AlertCount)
	}
}

func TestScoreIsAlertRatio(t *testing.T) {
	sc := New()
	sc.Record("tcp:9090", true)
	sc.Record("tcp:9090", true)
	sc.Record("tcp:9090", false)
	sc.Record("tcp:9090", false)
	e, _ := sc.Get("tcp:9090")
	want := 50.0
	if e.Score != want {
		t.Fatalf("expected score %.1f, got %.1f", want, e.Score)
	}
}

func TestGetMissingKeyReturnsFalse(t *testing.T) {
	sc := New()
	_, ok := sc.Get("tcp:1234")
	if ok {
		t.Fatal("expected missing key to return false")
	}
}

func TestSnapshotReturnsCopy(t *testing.T) {
	sc := New()
	sc.Record("tcp:80", true)
	sc.Record("udp:53", false)
	snap := sc.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap))
	}
}

func TestResetRemovesEntry(t *testing.T) {
	sc := New()
	sc.Record("tcp:443", false)
	sc.Reset("tcp:443")
	_, ok := sc.Get("tcp:443")
	if ok {
		t.Fatal("expected entry to be removed after reset")
	}
}

func TestIndependentKeysAreTrackedSeparately(t *testing.T) {
	sc := New()
	sc.Record("tcp:80", true)
	sc.Record("tcp:443", false)
	a, _ := sc.Get("tcp:80")
	b, _ := sc.Get("tcp:443")
	if a.AlertCount != 1 || b.AlertCount != 0 {
		t.Fatalf("entries should be independent: a=%+v b=%+v", a, b)
	}
}
