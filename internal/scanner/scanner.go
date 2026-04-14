package scanner

import (
	"fmt"
	"net"
	"time"
)

// PortState represents the state of a single port listener.
type PortState struct {
	Protocol string
	Address  string
	Port     int
	PID      int
	Process  string
	SeenAt   time.Time
}

// String returns a human-readable representation of a PortState.
func (p PortState) String() string {
	return fmt.Sprintf("%s %s:%d (pid=%d, process=%s)", p.Protocol, p.Address, p.Port, p.PID, p.Process)
}

// Scanner defines the interface for port scanning backends.
type Scanner interface {
	Scan() ([]PortState, error)
}

// IsLoopback returns true if the address is a loopback address.
func IsLoopback(addr string) bool {
	ip := net.ParseIP(addr)
	if ip == nil {
		return false
	}
	return ip.IsLoopback()
}

// DeduplicateStates removes duplicate PortState entries by protocol+address+port.
func DeduplicateStates(states []PortState) []PortState {
	seen := make(map[string]struct{})
	result := make([]PortState, 0, len(states))
	for _, s := range states {
		key := fmt.Sprintf("%s:%s:%d", s.Protocol, s.Address, s.Port)
		if _, ok := seen[key]; !ok {
			seen[key] = struct{}{}
			result = append(result, s)
		}
	}
	return result
}
