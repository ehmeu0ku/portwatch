package portping_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portping"
	"github.com/user/portwatch/internal/scanner"
)

func listenTCP(t *testing.T) (net.Listener, uint16) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })
	return ln, uint16(ln.Addr().(*net.TCPAddr).Port)
}

func TestFilterKeepsReachableStates(t *testing.T) {
	_, port := listenTCP(t)

	states := []scanner.PortState{
		{IP: "127.0.0.1", Port: port, Proto: "tcp"},
		{IP: "127.0.0.1", Port: 2, Proto: "tcp"}, // unreachable
	}

	f := portping.NewReachabilityFilter(500 * time.Millisecond)
	got := f.Filter(context.Background(), states)

	if len(got) != 1 {
		t.Fatalf("expected 1 reachable state, got %d", len(got))
	}
	if got[0].Port != port {
		t.Errorf("expected port %d, got %d", port, got[0].Port)
	}
}

func TestFilterAllUDPPassThrough(t *testing.T) {
	states := []scanner.PortState{
		{IP: "127.0.0.1", Port: 5000, Proto: "udp"},
		{IP: "127.0.0.1", Port: 5001, Proto: "udp"},
	}

	f := portping.NewReachabilityFilter(time.Second)
	got := f.Filter(context.Background(), states)

	if len(got) != 2 {
		t.Fatalf("expected 2 UDP states, got %d", len(got))
	}
}

func TestAnnotateReturnsMapForAllStates(t *testing.T) {
	_, port := listenTCP(t)

	states := []scanner.PortState{
		{IP: "127.0.0.1", Port: port, Proto: "tcp"},
	}

	f := portping.NewReachabilityFilter(time.Second)
	annotations := f.Annotate(context.Background(), states)

	if len(annotations) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(annotations))
	}
	for _, r := range annotations {
		if !r.Reachable {
			t.Error("expected reachable annotation for live port")
		}
	}
}

func TestFilterEmptyInputReturnsEmpty(t *testing.T) {
	f := portping.NewReachabilityFilter(time.Second)
	got := f.Filter(context.Background(), nil)
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %d elements", len(got))
	}
}
