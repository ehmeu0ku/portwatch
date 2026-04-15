package audit_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/portwatch/internal/audit"
)

func TestOpenCreatesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.jsonl")
	l, err := audit.Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer l.Close()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestRecordAndReadAll(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.jsonl")
	l, err := audit.Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}

	entries := []audit.Entry{
		{Kind: audit.EventNew, Proto: "tcp", Port: 8080, Addr: "0.0.0.0"},
		{Kind: audit.EventGone, Proto: "tcp", Port: 8080, Addr: "0.0.0.0"},
	}
	for _, e := range entries {
		if err := l.Record(e); err != nil {
			t.Fatalf("Record: %v", err)
		}
	}
	l.Close()

	got, err := audit.ReadAll(path)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if got[0].Kind != audit.EventNew {
		t.Errorf("entry 0 kind = %q, want NEW", got[0].Kind)
	}
	if got[1].Kind != audit.EventGone {
		t.Errorf("entry 1 kind = %q, want GONE", got[1].Kind)
	}
}

func TestRecordSetsTimestampWhenZero(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.jsonl")
	l, err := audit.Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer l.Close()

	before := time.Now().UTC()
	if err := l.Record(audit.Entry{Kind: audit.EventNew, Proto: "udp", Port: 53}); err != nil {
		t.Fatalf("Record: %v", err)
	}
	after := time.Now().UTC()
	l.Close()

	got, _ := audit.ReadAll(path)
	if got[0].Timestamp.Before(before) || got[0].Timestamp.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", got[0].Timestamp, before, after)
	}
}

func TestReadAllMissingFileReturnsError(t *testing.T) {
	_, err := audit.ReadAll("/nonexistent/path/audit.jsonl")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRecordPreservesExplicitTimestamp(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.jsonl")
	l, _ := audit.Open(path)
	defer l.Close()

	fixed := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	l.Record(audit.Entry{Kind: audit.EventNew, Proto: "tcp", Port: 443, Timestamp: fixed})
	l.Close()

	got, _ := audit.ReadAll(path)
	if !got[0].Timestamp.Equal(fixed) {
		t.Errorf("timestamp = %v, want %v", got[0].Timestamp, fixed)
	}
}
