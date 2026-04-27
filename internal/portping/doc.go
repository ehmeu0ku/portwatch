// Package portping provides active reachability probing for open ports.
//
// It complements passive scanner data by attempting real TCP connections
// (or UDP best-effort checks) so that portwatch can distinguish a port
// that is merely listed in /proc/net from one that is truly accepting
// traffic.
//
// Usage:
//
//	p := portping.New(2 * time.Second)
//	result := p.Probe(ctx, state)
//	if result.Reachable {
//		fmt.Println("port is live", result.Latency)
//	}
package portping
