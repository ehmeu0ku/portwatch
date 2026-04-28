// Package portfence enforces per-protocol port access fences.
// A fence defines a set of allowed ports for a given protocol;
// any port outside that set is considered a violation.
package portfence

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Violation is returned when a port falls outside its fence.
type Violation struct {
	Proto string
	Port  uint16
}

func (v Violation) Error() string {
	return fmt.Sprintf("port %d/%s is outside the allowed fence", v.Port, v.Proto)
}

// Fence holds allowed port sets keyed by protocol.
type Fence struct {
	mu      sync.RWMutex
	allowed map[string]map[uint16]struct{} // proto -> set of allowed ports
}

// New returns an empty Fence. Use Allow to populate it.
func New() *Fence {
	return &Fence{
		allowed: make(map[string]map[uint16]struct{}),
	}
}

// Allow registers port as permitted for the given protocol.
// proto should be "tcp" or "udp" (case-insensitive callers should normalise).
func (f *Fence) Allow(proto string, port uint16) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.allowed[proto]; !ok {
		f.allowed[proto] = make(map[uint16]struct{})
	}
	f.allowed[proto][port] = struct{}{}
}

// Check returns a Violation if the state's port is not in the fence for its
// protocol. If no fence has been defined for the protocol, Check returns nil
// (open by default).
func (f *Fence) Check(s scanner.PortState) error {
	f.mu.RLock()
	defer f.mu.RUnlock()
	set, ok := f.allowed[s.Proto]
	if !ok {
		return nil
	}
	if _, permitted := set[s.Port]; !permitted {
		return Violation{Proto: s.Proto, Port: s.Port}
	}
	return nil
}

// Filter returns only those states that pass the fence check.
func (f *Fence) Filter(states []scanner.PortState) []scanner.PortState {
	out := states[:0:0]
	for _, s := range states {
		if f.Check(s) == nil {
			out = append(out, s)
		}
	}
	return out
}

// Violations returns all states that violate the fence.
func (f *Fence) Violations(states []scanner.PortState) []Violation {
	var out []Violation
	for _, s := range states {
		if err := f.Check(s); err != nil {
			if v, ok := err.(Violation); ok {
				out = append(out, v)
			}
		}
	}
	return out
}
