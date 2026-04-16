package fingerprint_test

import (
	"testing"

	"github.com/user/portwatch/internal/fingerprint"
	"github.com/user/portwatch/internal/scanner"
)

func makeState(proto, ip string, port uint16) scanner.PortState {
	return scanner.PortState{Proto: proto, IP: ip, Port: port}
}

func TestOfProducesSameHashForSameState(t *testing.T) {
	s := makeState("tcp", "0.0.0.0", 8080)
	if fingerprint.Of(s) != fingerprint.Of(s) {
		t.Fatal("expected same fingerprint for identical state")
	}
}

func TestOfDiffersForDifferentPort(t *testing.T) {
	a := fingerprint.Of(makeState("tcp", "0.0.0.0", 8080))
	b := fingerprint.Of(makeState("tcp", "0.0.0.0", 9090))
	if a == b {
		t.Fatal("expected different fingerprints for different ports")
	}
}

func TestOfDiffersForDifferentProto(t *testing.T) {
	a := fingerprint.Of(makeState("tcp", "0.0.0.0", 80))
	b := fingerprint.Of(makeState("udp", "0.0.0.0", 80))
	if a == b {
		t.Fatal("expected different fingerprints for different protocols")
	}
}

func TestFingerprintStringIsHex(t *testing.T) {
	s := fingerprint.Of(makeState("tcp", "127.0.0.1", 22))
	if len(s.String()) != 16 {
		t.Fatalf("expected 16 hex chars, got %d", len(s.String()))
	}
}

func TestBuildCreatesMap(t *testing.T) {
	states := []scanner.PortState{
		makeState("tcp", "0.0.0.0", 80),
		makeState("tcp", "0.0.0.0", 443),
	}
	m := fingerprint.Build(states)
	if len(m) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(m))
	}
}

func TestAddedReturnsNewStates(t *testing.T) {
	prev := fingerprint.Build([]scanner.PortState{makeState("tcp", "0.0.0.0", 80)})
	next := fingerprint.Build([]scanner.PortState{
		makeState("tcp", "0.0.0.0", 80),
		makeState("tcp", "0.0.0.0", 9000),
	})
	added := fingerprint.Added(prev, next)
	if len(added) != 1 || added[0].Port != 9000 {
		t.Fatalf("expected port 9000 as added, got %v", added)
	}
}

func TestRemovedReturnsMissingStates(t *testing.T) {
	prev := fingerprint.Build([]scanner.PortState{
		makeState("tcp", "0.0.0.0", 80),
		makeState("tcp", "0.0.0.0", 3000),
	})
	next := fingerprint.Build([]scanner.PortState{makeState("tcp", "0.0.0.0", 80)})
	removed := fingerprint.Removed(prev, next)
	if len(removed) != 1 || removed[0].Port != 3000 {
		t.Fatalf("expected port 3000 as removed, got %v", removed)
	}
}
