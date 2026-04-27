package portdrift

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makeState(port uint16, proto, addr string) scanner.PortState {
	return scanner.PortState{Port: port, Proto: proto, Addr: addr}
}

func TestObserveFirstCallReturnsFalse(t *testing.T) {
	det := New()
	events := det.Observe([]scanner.PortState{makeState(80, "tcp", "0.0.0.0")})
	if len(events) != 0 {
		t.Fatalf("expected no events on first observe, got %d", len(events))
	}
}

func TestObserveSameAddrNoEvent(t *testing.T) {
	det := New()
	state := makeState(80, "tcp", "0.0.0.0")
	det.Observe([]scanner.PortState{state})
	events := det.Observe([]scanner.PortState{state})
	if len(events) != 0 {
		t.Fatalf("expected no events when address unchanged, got %d", len(events))
	}
}

func TestObserveAddrChangedReturnsDriftEvent(t *testing.T) {
	det := New()
	det.Observe([]scanner.PortState{makeState(443, "tcp", "127.0.0.1")})
	events := det.Observe([]scanner.PortState{makeState(443, "tcp", "0.0.0.0")})
	if len(events) != 1 {
		t.Fatalf("expected 1 drift event, got %d", len(events))
	}
	e := events[0]
	if e.Port != 443 || e.Proto != "tcp" {
		t.Errorf("unexpected port/proto: %d/%s", e.Port, e.Proto)
	}
	if e.OldAddr != "127.0.0.1" || e.NewAddr != "0.0.0.0" {
		t.Errorf("unexpected addrs: %s -> %s", e.OldAddr, e.NewAddr)
	}
}

func TestObserveMultiplePortsOnlyDriftedReported(t *testing.T) {
	det := New()
	det.Observe([]scanner.PortState{
		makeState(80, "tcp", "127.0.0.1"),
		makeState(443, "tcp", "0.0.0.0"),
	})
	events := det.Observe([]scanner.PortState{
		makeState(80, "tcp", "0.0.0.0"),  // drifted
		makeState(443, "tcp", "0.0.0.0"), // unchanged
	})
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Port != 80 {
		t.Errorf("expected port 80 to drift, got %d", events[0].Port)
	}
}

func TestForgetRemovesTracking(t *testing.T) {
	det := New()
	det.Observe([]scanner.PortState{makeState(8080, "tcp", "127.0.0.1")})
	det.Forget(8080, "tcp")
	// After forget, next observe should not report a drift
	events := det.Observe([]scanner.PortState{makeState(8080, "tcp", "0.0.0.0")})
	if len(events) != 0 {
		t.Fatalf("expected no events after forget, got %d", len(events))
	}
}

func TestLenReflectsTrackedPorts(t *testing.T) {
	det := New()
	if det.Len() != 0 {
		t.Fatalf("expected 0 initially, got %d", det.Len())
	}
	det.Observe([]scanner.PortState{
		makeState(80, "tcp", "0.0.0.0"),
		makeState(53, "udp", "0.0.0.0"),
	})
	if det.Len() != 2 {
		t.Fatalf("expected 2, got %d", det.Len())
	}
}

func TestDriftEventString(t *testing.T) {
	e := DriftEvent{Port: 80, Proto: "tcp", OldAddr: "127.0.0.1", NewAddr: "0.0.0.0"}
	s := e.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
	for _, want := range []string{"80", "tcp", "127.0.0.1", "0.0.0.0"} {
		if !containsStr(s, want) {
			t.Errorf("expected %q in %q", want, s)
		}
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && stringContains(s, sub))
}

func stringContains(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
