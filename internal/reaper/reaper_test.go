package reaper

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeFile(t *testing.T, dir, name string, age time.Duration) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	mod := time.Now().Add(-age)
	if err := os.Chtimes(p, mod, mod); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestPruneRemovesOldFiles(t *testing.T) {
	dir := t.TempDir()
	old := writeFile(t, dir, "snap-old.json", 2*time.Hour)
	_ = writeFile(t, dir, "snap-new.json", 10*time.Second)

	r := New(dir, time.Hour, time.Minute, "snap-*.json", nil)
	n := r.Prune()

	if n != 1 {
		t.Fatalf("expected 1 removal, got %d", n)
	}
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Error("old file should have been removed")
	}
}

func TestPruneKeepsRecentFiles(t *testing.T) {
	dir := t.TempDir()
	_ = writeFile(t, dir, "snap-recent.json", 30*time.Second)

	r := New(dir, time.Hour, time.Minute, "snap-*.json", nil)
	n := r.Prune()

	if n != 0 {
		t.Fatalf("expected 0 removals, got %d", n)
	}
}

func TestPruneEmptyDirReturnsZero(t *testing.T) {
	dir := t.TempDir()
	r := New(dir, time.Hour, time.Minute, "*.json", nil)
	if n := r.Prune(); n != 0 {
		t.Fatalf("expected 0, got %d", n)
	}
}

func TestPruneGlobMismatchSkips(t *testing.T) {
	dir := t.TempDir()
	_ = writeFile(t, dir, "audit.log", 5*time.Hour)

	r := New(dir, time.Hour, time.Minute, "snap-*.json", nil)
	n := r.Prune()

	if n != 0 {
		t.Fatalf("glob should not match .log file, got %d removals", n)
	}
}
