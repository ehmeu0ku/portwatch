package enricher_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/enricher"
	"github.com/user/portwatch/internal/process"
	"github.com/user/portwatch/internal/scanner"
)

// buildFakeProc creates a temporary /proc-like directory tree that maps
// the given inode to a process, reusing the helper from the process package
// tests via a local reimplementation.
func makeState(port uint16, inode uint64) scanner.PortState {
	return scanner.PortState{
		Port:  port,
		Proto: "tcp",
		Addr:  "0.0.0.0",
		Inode: inode,
	}
}

func TestEnrichNilProcessWhenInodeUnknown(t *testing.T) {
	// Use a real resolver pointed at a non-existent proc root so every
	// lookup fails gracefully.
	r := process.NewResolver("/proc/nonexistent/fake")
	e := enricher.New(r)

	states := []scanner.PortState{makeState(8080, 99999)}
	result := e.Enrich(states)

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if result[0].Process != nil {
		t.Errorf("expected nil Process for unknown inode, got %v", result[0].Process)
	}
}

func TestEnrichPreservesPortState(t *testing.T) {
	r := process.NewResolver("/proc/nonexistent/fake")
	e := enricher.New(r)

	s := makeState(443, 12345)
	result := e.Enrich([]scanner.PortState{s})

	if result[0].Port != 443 {
		t.Errorf("expected port 443, got %d", result[0].Port)
	}
	if result[0].Proto != "tcp" {
		t.Errorf("expected proto tcp, got %s", result[0].Proto)
	}
}

func TestEnrichOneReturnsEnrichedState(t *testing.T) {
	r := process.NewResolver("/proc/nonexistent/fake")
	e := enricher.New(r)

	s := makeState(22, 7777)
	es := e.EnrichOne(s)

	if es.Port != 22 {
		t.Errorf("expected port 22, got %d", es.Port)
	}
	if es.Process != nil {
		t.Errorf("expected nil Process, got %v", es.Process)
	}
}

func TestEnrichedStateStringWithoutProcess(t *testing.T) {
	r := process.NewResolver("/proc/nonexistent/fake")
	e := enricher.New(r)

	es := e.EnrichOne(makeState(80, 0))
	got := es.String()

	if strings.Contains(got, "[") {
		t.Errorf("expected no process bracket in string, got: %s", got)
	}
}

func TestEnrichEmptySlice(t *testing.T) {
	r := process.NewResolver("/proc/nonexistent/fake")
	e := enricher.New(r)

	result := e.Enrich([]scanner.PortState{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d items", len(result))
	}
}
