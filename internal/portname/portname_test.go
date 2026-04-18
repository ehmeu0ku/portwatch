package portname_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/portname"
)

func TestLookupWellKnownPort(t *testing.T) {
	r := portname.New(nil)
	if got := r.Lookup(22); got != "ssh" {
		t.Fatalf("expected ssh, got %s", got)
	}
}

func TestLookupHTTPS(t *testing.T) {
	r := portname.New(nil)
	if got := r.Lookup(443); got != "https" {
		t.Fatalf("expected https, got %s", got)
	}
}

func TestLookupUnknownReturnsFallback(t *testing.T) {
	r := portname.New(nil)
	if got := r.Lookup(9999); got != "port-9999" {
		t.Fatalf("expected port-9999, got %s", got)
	}
}

func TestCustomMappingOverridesBuiltin(t *testing.T) {
	r := portname.New(map[uint16]string{80: "my-http"})
	if got := r.Lookup(80); got != "my-http" {
		t.Fatalf("expected my-http, got %s", got)
	}
}

func TestCustomMappingAddsNewPort(t *testing.T) {
	r := portname.New(map[uint16]string{12345: "my-service"})
	if got := r.Lookup(12345); got != "my-service" {
		t.Fatalf("expected my-service, got %s", got)
	}
}

func TestKnownReturnsTrueForWellKnown(t *testing.T) {
	r := portname.New(nil)
	if !r.Known(22) {
		t.Fatal("expected port 22 to be known")
	}
}

func TestKnownReturnsFalseForUnknown(t *testing.T) {
	r := portname.New(nil)
	if r.Known(9999) {
		t.Fatal("expected port 9999 to be unknown")
	}
}

func TestKnownReturnsTrueForCustomPort(t *testing.T) {
	r := portname.New(map[uint16]string{12345: "custom"})
	if !r.Known(12345) {
		t.Fatal("expected custom port to be known")
	}
}
