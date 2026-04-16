package history

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makeState(port uint16, proto string) scanner.PortState {
	return scanner.PortState{Port: port, Proto: proto, LocalAddr: "0.0.0.0"}
}

func TestNewHistoryIsEmpty(t *testing.T) {
	h := New(10)
	if len(h.Entries()) != 0 {
		t.Fatalf("expected empty history, got %d entries", len(h.Entries()))
	}
}

func TestRecordAddsEntry(t *testing.T) {
	h := New(10)
	h.Record(makeState(8080, "tcp"), "new")
	entries := h.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Event != "new" {
		t.Errorf("expected event 'new', got %q", entries[0].Event)
	}
	if entries[0].State.Port != 8080 {
		t.Errorf("expected port 8080, got %d", entries[0].State.Port)
	}
	if entries[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestHistoryEvictsOldestWhenFull(t *testing.T) {
	h := New(3)
	for i := uint16(1); i <= 4; i++ {
		h.Record(makeState(i, "tcp"), "new")
	}
	entries := h.Entries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", len(entries))
	}
	if entries[0].State.Port != 2 {
		t.Errorf("expected oldest entry to be port 2, got %d", entries[0].State.Port)
	}
}

func TestEntriesReturnsCopy(t *testing.T) {
	h := New(10)
	h.Record(makeState(443, "tcp"), "new")
	e1 := h.Entries()
	e1[0].Event = "mutated"
	e2 := h.Entries()
	if e2[0].Event == "mutated" {
		t.Error("Entries() should return a copy, not a reference")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	h1 := New(10)
	h1.Record(makeState(22, "tcp"), "new")
	h1.Record(makeState(22, "tcp"), "gone")
	if err := h1.Save(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	h2 := New(10)
	if err := h2.Load(path); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	entries := h2.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries after load, got %d", len(entries))
	}
	if entries[0].State.Port != 22 || entries[0].Event != "new" {
		t.Errorf("unexpected first entry: port=%d event=%q", entries[0].State.Port, entries[0].Event)
	}
	if entries[1].State.Port != 22 || entries[1].Event != "gone" {
		t.Errorf("unexpected second entry: port=%d event=%q", entries[1].State.Port, entries[1].Event)
	}
}

func TestLoadNonExistentFileIsNoop(t *testing.T) {
	h := New(10)
	if err := h.Load("/nonexistent/path/history.json"); err != nil {
		t.Errorf("expected no error for missing file, got %v", err)
	}
}

func TestDefaultMaxSize(t *testing.T) {
	h := New(0)
	if h.maxSize != 500 {
		t.Errorf("expected default maxSize 500, got %d", h.maxSize)
	}
}

func TestSaveCreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.json")
	h := New(10)
	h.Record(makeState(80, "tcp"), "new")
	if err := h.Save(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected file to exist after Save")
	}
}
