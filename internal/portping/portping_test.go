package portping_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portping"
	"github.com/user/portwatch/internal/scanner"
)

func makeState(ip string, port uint16, proto string) scanner.PortState {
	return scanner.PortState{IP: ip, Port: port, Proto: proto}
}

func TestProbeReachableTCPPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	addr := ln.Addr().(*net.TCPAddr)
	s := makeState("127.0.0.1", uint16(addr.Port), "tcp")

	p := portping.New(2 * time.Second)
	res := p.Probe(context.Background(), s)

	if !res.Reachable {
		t.Fatalf("expected reachable, got err: %v", res.Err)
	}
	if res.Latency <= 0 {
		t.Error("expected positive latency")
	}
}

func TestProbeUnreachableTCPPort(t *testing.T) {
	// Port 1 is almost certainly closed.
	s := makeState("127.0.0.1", 1, "tcp")
	p := portping.New(300 * time.Millisecond)
	res := p.Probe(context.Background(), s)

	if res.Reachable {
		t.Fatal("expected unreachable")
	}
	if res.Err == nil {
		t.Error("expected non-nil error")
	}
}

func TestProbeUDPAlwaysReachable(t *testing.T) {
	s := makeState("127.0.0.1", 9999, "udp")
	p := portping.New(time.Second)
	res := p.Probe(context.Background(), s)

	if !res.Reachable {
		t.Error("UDP probe should always return reachable")
	}
	if res.Err != nil {
		t.Errorf("unexpected error: %v", res.Err)
	}
}

func TestProbeAllReturnsOneResultPerState(t *testing.T) {
	states := []scanner.PortState{
		makeState("127.0.0.1", 9000, "udp"),
		makeState("127.0.0.1", 9001, "udp"),
	}
	p := portping.New(time.Second)
	results := p.ProbeAll(context.Background(), states)

	if len(results) != len(states) {
		t.Fatalf("expected %d results, got %d", len(states), len(results))
	}
}

func TestNewDefaultsZeroTimeout(t *testing.T) {
	// Should not panic or hang; zero timeout is replaced with 2s internally.
	p := portping.New(0)
	if p == nil {
		t.Fatal("expected non-nil prober")
	}
}

func TestProbeContextCancellation(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	addr := ln.Addr().(*net.TCPAddr)
	s := makeState("127.0.0.1", uint16(addr.Port), "tcp")

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	p := portping.New(2 * time.Second)
	res := p.Probe(ctx, s)
	// With a cancelled context the dial may or may not succeed depending on
	// timing; we only assert no panic occurs.
	_ = res
}
