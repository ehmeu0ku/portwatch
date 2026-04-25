package portmatch_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/portmatch"
	"github.com/user/portwatch/internal/scanner"
)

func writeJSON(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "exprs.json")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoadFileValid(t *testing.T) {
	p := writeJSON(t, `[
		{"port": 80, "proto": "tcp"},
		{"port_min": 8000, "port_max": 9000, "tag": "dev"}
	]`)
	m, err := portmatch.LoadFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.AnyMatch(scanner.PortState{Port: 80, Proto: "tcp"}) {
		t.Error("expected match on port 80/tcp")
	}
	if !m.AnyMatch(scanner.PortState{Port: 8500, Proto: "udp", Tags: []string{"dev"}}) {
		t.Error("expected match in range with tag dev")
	}
}

func TestLoadFileMissing(t *testing.T) {
	_, err := portmatch.LoadFile("/nonexistent/path/exprs.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadFileUnknownProto(t *testing.T) {
	p := writeJSON(t, `[{"port": 443, "proto": "sctp"}]`)
	_, err := portmatch.LoadFile(p)
	if err == nil {
		t.Fatal("expected error for unknown proto")
	}
}

func TestLoadFileBadRange(t *testing.T) {
	p := writeJSON(t, `[{"port_min": 9000, "port_max": 8000}]`)
	_, err := portmatch.LoadFile(p)
	if err == nil {
		t.Fatal("expected error for inverted range")
	}
}

func TestLoadFileEmptyArray(t *testing.T) {
	p := writeJSON(t, `[]`)
	m, err := portmatch.LoadFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.AnyMatch(scanner.PortState{Port: 80, Proto: "tcp"}) {
		t.Error("empty matcher should not match anything")
	}
}
