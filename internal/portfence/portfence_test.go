package portfence_test

import (
	"testing"

	"github.com/user/portwatch/internal/portfence"
	"github.com/user/portwatch/internal/scanner"
)

func makeState(proto string, port uint16) scanner.PortState {
	return scanner.PortState{Proto: proto, Port: port, IP: "0.0.0.0"}
}

func TestCheckAllowedPortReturnsNil(t *testing.T) {
	f := portfence.New()
	f.Allow("tcp", 443)
	if err := f.Check(makeState("tcp", 443)); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCheckForbiddenPortReturnsViolation(t *testing.T) {
	f := portfence.New()
	f.Allow("tcp", 80)
	err := f.Check(makeState("tcp", 9000))
	if err == nil {
		t.Fatal("expected violation, got nil")
	}
	v, ok := err.(portfence.Violation)
	if !ok {
		t.Fatalf("expected Violation type, got %T", err)
	}
	if v.Port != 9000 || v.Proto != "tcp" {
		t.Fatalf("unexpected violation fields: %+v", v)
	}
}

func TestCheckNoFenceDefinedReturnsNil(t *testing.T) {
	f := portfence.New()
	// no fence for udp
	if err := f.Check(makeState("udp", 53)); err != nil {
		t.Fatalf("expected nil for unfenced proto, got %v", err)
	}
}

func TestFilterKeepsAllowedStates(t *testing.T) {
	f := portfence.New()
	f.Allow("tcp", 80)
	f.Allow("tcp", 443)
	states := []scanner.PortState{
		makeState("tcp", 80),
		makeState("tcp", 8080),
		makeState("tcp", 443),
	}
	got := f.Filter(states)
	if len(got) != 2 {
		t.Fatalf("expected 2 states, got %d", len(got))
	}
}

func TestViolationsReturnsOnlyForbidden(t *testing.T) {
	f := portfence.New()
	f.Allow("tcp", 22)
	states := []scanner.PortState{
		makeState("tcp", 22),
		makeState("tcp", 4444),
		makeState("tcp", 31337),
	}
	vs := f.Violations(states)
	if len(vs) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(vs))
	}
}

func TestViolationErrorString(t *testing.T) {
	v := portfence.Violation{Proto: "tcp", Port: 1234}
	got := v.Error()
	if got == "" {
		t.Fatal("expected non-empty error string")
	}
}

func TestAllowMultipleProtocols(t *testing.T) {
	f := portfence.New()
	f.Allow("tcp", 80)
	f.Allow("udp", 53)
	if err := f.Check(makeState("tcp", 80)); err != nil {
		t.Fatalf("tcp 80 should be allowed: %v", err)
	}
	if err := f.Check(makeState("udp", 53)); err != nil {
		t.Fatalf("udp 53 should be allowed: %v", err)
	}
	if err := f.Check(makeState("tcp", 53)); err == nil {
		t.Fatal("tcp 53 should be a violation")
	}
}
