package statestore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/statestore"
)

func makeState(port uint16, proto string) scanner.PortState {
	return scanner.PortState{Port: port, Proto: proto, State: "LISTEN"}
}

func TestNewStoreIsEmpty(t *testing.T) {
	s := statestore.New(filepath.Join(t.TempDir(), "state.json"))
	if got := s.Get(); len(got) != 0 {
		t.Fatalf("expected empty, got %d entries", len(got))
	}
}

func TestSetAndGet(t *testing.T) {
	dir := t.TempDir()
	s := statestore.New(filepath.Join(dir, "state.json"))
	states := []scanner.PortState{makeState(80, "tcp"), makeState(443, "tcp")}
	if err := s.Set(states); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got := s.Get()
	if len(got) != 2 {
		t.Fatalf("expected 2 states, got %d", len(got))
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s := statestore.New(path)
	_ = s.Set([]scanner.PortState{makeState(22, "tcp")})

	loaded, err := statestore.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	got := loaded.Get()
	if len(got) != 1 || got[0].Port != 22 {
		t.Fatalf("unexpected states: %+v", got)
	}
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	s, err := statestore.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(s.Get()) != 0 {
		t.Fatal("expected empty store for missing file")
	}
}

func TestUpdatedAtSetAfterSet(t *testing.T) {
	s := statestore.New(filepath.Join(t.TempDir(), "state.json"))
	if !s.UpdatedAt().IsZero() {
		t.Fatal("expected zero time before Set")
	}
	_ = s.Set([]scanner.PortState{makeState(8080, "tcp")})
	if s.UpdatedAt().IsZero() {
		t.Fatal("expected non-zero time after Set")
	}
}

func TestGetReturnsCopy(t *testing.T) {
	s := statestore.New(filepath.Join(t.TempDir(), "state.json"))
	_ = s.Set([]scanner.PortState{makeState(9000, "tcp")})
	a := s.Get()
	a[0].Port = 1
	if s.Get()[0].Port != 9000 {
		t.Fatal("Get should return a copy, not a reference")
	}
}

func TestSetWritesFileToDisk(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s := statestore.New(path)
	_ = s.Set([]scanner.PortState{makeState(3306, "tcp")})
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file on disk: %v", err)
	}
}
