package retention_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/retention"
)

type record struct{ ts time.Time }

func (r record) Timestamp() time.Time { return r.ts }

func TestAddAndLen(t *testing.T) {
	s := retention.New(time.Minute)
	s.Add(record{ts: time.Now()})
	s.Add(record{ts: time.Now()})
	if s.Len() != 2 {
		t.Fatalf("expected 2, got %d", s.Len())
	}
}

func TestEntriesReturnsCopy(t *testing.T) {
	s := retention.New(time.Minute)
	s.Add(record{ts: time.Now()})
	e1 := s.Entries()
	e1[0] = record{ts: time.Time{}}
	e2 := s.Entries()
	if e2[0].Timestamp().IsZero() {
		t.Fatal("copy was not independent")
	}
}

func TestPrunesExpiredEntries(t *testing.T) {
	s := retention.New(time.Minute)
	old := time.Now().Add(-2 * time.Minute)
	s.Add(record{ts: old})
	s.Add(record{ts: time.Now()})
	if s.Len() != 1 {
		t.Fatalf("expected 1 live entry, got %d", s.Len())
	}
}

func TestAllExpiredReturnsEmpty(t *testing.T) {
	s := retention.New(time.Minute)
	s.Add(record{ts: time.Now().Add(-5 * time.Minute)})
	if s.Len() != 0 {
		t.Fatal("expected empty store after full expiry")
	}
}

func TestLenPrunesBeforeCounting(t *testing.T) {
	s := retention.New(30 * time.Second)
	s.Add(record{ts: time.Now().Add(-60 * time.Second)})
	s.Add(record{ts: time.Now()})
	if got := s.Len(); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}
