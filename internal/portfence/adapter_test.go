package portfence_test

import (
	"testing"

	"github.com/example/portwatch/internal/portfence"
	"github.com/example/portwatch/internal/scanner"
)

func makeAdapterState(port uint16, proto string) scanner.PortState {
	return scanner.PortState{
		Port:  port,
		Proto: proto,
		Addr:  "0.0.0.0",
	}
}

func TestAdapterFilterPassesAllowedStates(t *testing.T) {
	fence := portfence.New(portfence.Config{
		DefaultAllow: false,
	})
	fence.Allow(80, "tcp")
	fence.Allow(443, "tcp")

	states := []scanner.PortState{
		makeAdapterState(80, "tcp"),
		makeAdapterState(443, "tcp"),
		makeAdapterState(9999, "tcp"),
	}

	filter := portfence.NewFilter(fence)
	result := filter.Apply(states)

	if len(result) != 2 {
		t.Fatalf("expected 2 allowed states, got %d", len(result))
	}
	for _, s := range result {
		if s.Port == 9999 {
			t.Errorf("forbidden port 9999 should not appear in result")
		}
	}
}

func TestAdapterFilterBlocksForbiddenStates(t *testing.T) {
	fence := portfence.New(portfence.Config{
		DefaultAllow: true,
	})
	fence.Deny(22, "tcp")

	states := []scanner.PortState{
		makeAdapterState(22, "tcp"),
		makeAdapterState(80, "tcp"),
	}

	filter := portfence.NewFilter(fence)
	result := filter.Apply(states)

	if len(result) != 1 {
		t.Fatalf("expected 1 state, got %d", len(result))
	}
	if result[0].Port != 80 {
		t.Errorf("expected port 80 to pass, got %d", result[0].Port)
	}
}

func TestAdapterFilterEmptyInputReturnsEmpty(t *testing.T) {
	fence := portfence.New(portfence.Config{DefaultAllow: true})
	filter := portfence.NewFilter(fence)
	result := filter.Apply(nil)
	if len(result) != 0 {
		t.Errorf("expected empty result for nil input, got %d", len(result))
	}
}

func TestAdapterFilterDefaultAllowPassesUnknown(t *testing.T) {
	fence := portfence.New(portfence.Config{DefaultAllow: true})
	states := []scanner.PortState{
		makeAdapterState(12345, "tcp"),
	}
	filter := portfence.NewFilter(fence)
	result := filter.Apply(states)
	if len(result) != 1 {
		t.Errorf("expected 1 state with default-allow, got %d", len(result))
	}
}

func TestAdapterFilterDefaultDenyBlocksUnknown(t *testing.T) {
	fence := portfence.New(portfence.Config{DefaultAllow: false})
	states := []scanner.PortState{
		makeAdapterState(12345, "tcp"),
	}
	filter := portfence.NewFilter(fence)
	result := filter.Apply(states)
	if len(result) != 0 {
		t.Errorf("expected 0 states with default-deny, got %d", len(result))
	}
}
