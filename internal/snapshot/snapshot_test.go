package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func makeState(proto, addr string, port uint16) scanner.PortState {
	return scanner.PortState{Proto: proto, LocalAddr: addr, LocalPort: port}
}

func TestNewSnapshotSetsTimestamp(t *testing.T) {
	before := time.Now().UTC()
	snap := snapshot.New(nil)
	after := time.Now().UTC()

	if snap.CapturedAt.Before(before) || snap.CapturedAt.After(after) {
		t.Errorf("CapturedAt %v not in expected range [%v, %v]", snap.CapturedAt, before, after)
	}
}

func TestNewSnapshotStoresStates(t *testing.T) {
	states := []scanner.PortState{
		makeState("tcp", "0.0.0.0", 80),
		makeState("tcp", "0.0.0.0", 443),
	}
	snap := snapshot.New(states)
	if len(snap.States) != 2 {
		t.Fatalf("expected 2 states, got %d", len(snap.States))
	}
}

func TestSaveAndLoad(t *testing.T) {
	states := []scanner.PortState{
		makeState("tcp", "127.0.0.1", 8080),
	}
	snap := snapshot.New(states)

	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	if err := snap.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(loaded.States) != 1 {
		t.Fatalf("expected 1 state after load, got %d", len(loaded.States))
	}
	if loaded.States[0].LocalPort != 8080 {
		t.Errorf("expected port 8080, got %d", loaded.States[0].LocalPort)
	}
	if loaded.CapturedAt.IsZero() {
		t.Error("CapturedAt should not be zero after load")
	}
}

func TestLoadMissingFileReturnsError(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error loading missing file, got nil")
	}
}

func TestSaveInvalidPathReturnsError(t *testing.T) {
	snap := snapshot.New(nil)
	err := snap.Save("/nonexistent/dir/snap.json")
	if err == nil {
		t.Error("expected error saving to invalid path, got nil")
	}
}

func TestSummaryFormat(t *testing.T) {
	states := []scanner.PortState{
		makeState("udp", "0.0.0.0", 53),
		makeState("udp", "0.0.0.0", 123),
	}
	snap := snapshot.New(states)
	summary := snap.Summary()
	if summary == "" {
		t.Error("Summary should not be empty")
	}
	_ = os.Stdout // ensure os import used via t.TempDir indirectly
}
