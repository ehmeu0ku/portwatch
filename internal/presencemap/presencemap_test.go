package presencemap_test

import (
	"testing"

	"github.com/user/portwatch/internal/presencemap"
)

func TestObserveReturnsFalseBeforeThreshold(t *testing.T) {
	p := presencemap.New(3)
	if p.Observe("tcp:8080") {
		t.Fatal("expected false on first observation")
	}
	if p.Observe("tcp:8080") {
		t.Fatal("expected false on second observation")
	}
}

func TestObserveReturnsTrueAtThreshold(t *testing.T) {
	p := presencemap.New(3)
	p.Observe("tcp:8080")
	p.Observe("tcp:8080")
	if !p.Observe("tcp:8080") {
		t.Fatal("expected true at threshold")
	}
}

func TestObserveReturnsTrueBeyondThreshold(t *testing.T) {
	p := presencemap.New(2)
	p.Observe("tcp:443")
	p.Observe("tcp:443")
	if !p.Observe("tcp:443") {
		t.Fatal("expected true beyond threshold")
	}
}

func TestForgetResetsCounter(t *testing.T) {
	p := presencemap.New(2)
	p.Observe("tcp:9000")
	p.Observe("tcp:9000")
	p.Forget("tcp:9000")
	if p.Count("tcp:9000") != 0 {
		t.Fatalf("expected 0 after forget, got %d", p.Count("tcp:9000"))
	}
	if p.Observe("tcp:9000") {
		t.Fatal("expected false after forget and single observe with threshold 2")
	}
}

func TestStableReflectsThreshold(t *testing.T) {
	p := presencemap.New(2)
	if p.Stable("tcp:22") {
		t.Fatal("expected not stable before any observation")
	}
	p.Observe("tcp:22")
	if p.Stable("tcp:22") {
		t.Fatal("expected not stable after one observation")
	}
	p.Observe("tcp:22")
	if !p.Stable("tcp:22") {
		t.Fatal("expected stable after reaching threshold")
	}
}

func TestLenTracksKeys(t *testing.T) {
	p := presencemap.New(1)
	if p.Len() != 0 {
		t.Fatal("expected empty map")
	}
	p.Observe("tcp:80")
	p.Observe("udp:53")
	if p.Len() != 2 {
		t.Fatalf("expected 2, got %d", p.Len())
	}
	p.Forget("tcp:80")
	if p.Len() != 1 {
		t.Fatalf("expected 1 after forget, got %d", p.Len())
	}
}

func TestThresholdBelowOneDefaultsToOne(t *testing.T) {
	p := presencemap.New(0)
	if !p.Observe("tcp:1234") {
		t.Fatal("threshold 0 should default to 1, first observe should return true")
	}
}
