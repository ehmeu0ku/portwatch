package monitor

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func ps(proto, addr string, port uint16) scanner.PortState {
	return scanner.PortState{Proto: proto, Addr: addr, Port: port}
}

func TestDiffNoPrevious(t *testing.T) {
	current := []scanner.PortState{ps("tcp", "0.0.0.0", 8080)}
	added, removed := diff(nil, current)
	if len(added) != 1 || added[0].Port != 8080 {
		t.Fatalf("expected 1 added port, got %v", added)
	}
	if len(removed) != 0 {
		t.Fatalf("expected 0 removed ports, got %v", removed)
	}
}

func TestDiffGonePort(t *testing.T) {
	prev := []scanner.PortState{ps("tcp", "0.0.0.0", 9000)}
	added, removed := diff(prev, nil)
	if len(added) != 0 {
		t.Fatalf("expected 0 added ports, got %v", added)
	}
	if len(removed) != 1 || removed[0].Port != 9000 {
		t.Fatalf("expected 1 removed port, got %v", removed)
	}
}

func TestDiffNoChange(t *testing.T) {
	states := []scanner.PortState{
		ps("tcp", "0.0.0.0", 80),
		ps("udp", "0.0.0.0", 53),
	}
	added, removed := diff(states, states)
	if len(added) != 0 || len(removed) != 0 {
		t.Fatalf("expected no changes, got added=%v removed=%v", added, removed)
	}
}

func TestDiffMixed(t *testing.T) {
	prev := []scanner.PortState{
		ps("tcp", "0.0.0.0", 80),
		ps("tcp", "0.0.0.0", 443),
	}
	curr := []scanner.PortState{
		ps("tcp", "0.0.0.0", 443),
		ps("tcp", "0.0.0.0", 8080),
	}
	added, removed := diff(prev, curr)
	if len(added) != 1 || added[0].Port != 8080 {
		t.Fatalf("expected port 8080 added, got %v", added)
	}
	if len(removed) != 1 || removed[0].Port != 80 {
		t.Fatalf("expected port 80 removed, got %v", removed)
	}
}
