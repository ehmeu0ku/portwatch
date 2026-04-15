package tagger_test

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/tagger"
)

func makeState(port uint16) scanner.PortState {
	return scanner.PortState{
		Proto:   "tcp",
		Port:    port,
		Address: "0.0.0.0",
	}
}

func TestTagKnownPort(t *testing.T) {
	tr := tagger.New(nil)
	got := tr.Tag(makeState(80))
	if got != "http" {
		t.Fatalf("expected http, got %q", got)
	}
}

func TestTagUnknownPortReturnsEmpty(t *testing.T) {
	tr := tagger.New(nil)
	got := tr.Tag(makeState(9999))
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestCustomMappingOverridesWellKnown(t *testing.T) {
	tr := tagger.New(map[uint16]string{80: "my-proxy"})
	got := tr.Tag(makeState(80))
	if got != "my-proxy" {
		t.Fatalf("expected my-proxy, got %q", got)
	}
}

func TestCustomMappingAddsNewPort(t *testing.T) {
	tr := tagger.New(map[uint16]string{9200: "elasticsearch"})
	got := tr.Tag(makeState(9200))
	if got != "elasticsearch" {
		t.Fatalf("expected elasticsearch, got %q", got)
	}
}

func TestTagAllReturnsOnlyTaggedPorts(t *testing.T) {
	tr := tagger.New(nil)
	states := []scanner.PortState{
		makeState(22),
		makeState(9999),
		makeState(443),
	}
	tags := tr.TagAll(states)

	if len(tags) != 2 {
		t.Fatalf("expected 2 tagged ports, got %d", len(tags))
	}
	if tags[22] != "ssh" {
		t.Errorf("expected ssh for port 22, got %q", tags[22])
	}
	if tags[443] != "https" {
		t.Errorf("expected https for port 443, got %q", tags[443])
	}
	if _, ok := tags[9999]; ok {
		t.Error("port 9999 should not be tagged")
	}
}

func TestTagAllEmptyInput(t *testing.T) {
	tr := tagger.New(nil)
	tags := tr.TagAll(nil)
	if len(tags) != 0 {
		t.Fatalf("expected empty map, got %v", tags)
	}
}
