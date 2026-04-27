// Package portping probes whether a TCP/UDP port is actively accepting
// connections, adding a reachability dimension beyond mere listener presence.
package portping

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Result holds the outcome of a single probe attempt.
type Result struct {
	State     scanner.PortState
	Reachable bool
	Latency   time.Duration
	Err       error
}

// Prober probes ports for reachability.
type Prober struct {
	timeout time.Duration
}

// New creates a Prober with the given per-probe timeout.
func New(timeout time.Duration) *Prober {
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	return &Prober{timeout: timeout}
}

// Probe attempts a TCP dial to the port described by state.
// UDP ports are marked reachable without a dial (best-effort).
func (p *Prober) Probe(ctx context.Context, state scanner.PortState) Result {
	res := Result{State: state}

	if state.Proto == "udp" {
		// UDP has no handshake; mark as reachable by convention.
		res.Reachable = true
		return res
	}

	addr := fmt.Sprintf("%s:%d", state.IP, state.Port)
	start := time.Now()

	dialCtx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var d net.Dialer
	conn, err := d.DialContext(dialCtx, "tcp", addr)
	res.Latency = time.Since(start)

	if err != nil {
		res.Err = err
		return res
	}
	_ = conn.Close()
	res.Reachable = true
	return res
}

// ProbeAll probes every state in the slice and returns corresponding results.
func (p *Prober) ProbeAll(ctx context.Context, states []scanner.PortState) []Result {
	out := make([]Result, len(states))
	for i, s := range states {
		out[i] = p.Probe(ctx, s)
	}
	return out
}
