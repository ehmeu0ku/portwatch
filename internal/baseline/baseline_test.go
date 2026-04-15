package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/scanner"
)

func makeState(proto, addr string, port uint16) scanner.PortState {
	return scanner.PortState{Proto: proto, Address: addr, Port: port}
}

func TestNewBaselineIsEmpty(t *testing.T) {
	b := baseline.New("/tmp/unused.json")
	if len(b.Entries()) != 0 {
		t.Fatal("expected empty baseline")
	}
}

func TestAddAndContains(t *testing.T) {
	b := baseline.New("/tmp/unused.json")
	s := makeState("tcp", "0.0.0.0", 8080)
	if b.Contains(s) {
		t.Fatal("should not contain state before Add")
	}
	b.Add(s)
	if !b.Contains(s) {
		t.Fatal("should contain state after Add")
	}
}

func TestContainsDifferentPort(t *testing.T) {
	b := baseline.New("/tmp/unused.json")
	b.Add(makeState("tcp", "0.0.0.0", 8080))
	if b.Contains(makeState("tcp", "0.0.0.0", 9090)) {
		t.Fatal("should not contain a different port")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	b := baseline.New(path)
	b.Add(makeState("tcp", "0.0.0.0", 22))
	b.Add(makeState("tcp", "0.0.0.0", 443))

	if err := b.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !loaded.Contains(makeState("tcp", "0.0.0.0", 22)) {
		t.Error("loaded baseline missing port 22")
	}
	if !loaded.Contains(makeState("tcp", "0.0.0.0", 443)) {
		t.Error("loaded baseline missing port 443")
	}
	if len(loaded.Entries()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(loaded.Entries()))
	}
}

func TestLoadMissingFile(t *testing.T) {
	b, err := baseline.Load("/nonexistent/path/baseline.json")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil baseline")
	}
}

func TestLoadCorruptFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0o600)
	_, err := baseline.Load(path)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}
