// Package portpolicy evaluates whether a listening port is permitted
// according to a set of declarative allow/deny rules.
package portpolicy

import (
	"fmt"
	"sync"
)

// Action describes whether a rule permits or denies a port.
type Action int

const (
	Allow Action = iota
	Deny
)

func (a Action) String() string {
	if a == Allow {
		return "allow"
	}
	return "deny"
}

// Rule matches a port/protocol pair and assigns an action.
type Rule struct {
	Port     uint16
	Proto    string // "tcp" or "udp", empty matches both
	Action   Action
	Comment  string
}

// Policy holds an ordered list of rules. The first matching rule wins.
// If no rule matches the default action is Deny.
type Policy struct {
	mu      sync.RWMutex
	rules   []Rule
	default_ Action
}

// New returns a Policy with the given default action.
func New(defaultAction Action) *Policy {
	return &Policy{default_: defaultAction}
}

// Add appends a rule to the policy.
func (p *Policy) Add(r Rule) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.rules = append(p.rules, r)
}

// Evaluate returns the action that applies to the given port and protocol.
func (p *Policy) Evaluate(port uint16, proto string) Action {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, r := range p.rules {
		if r.Port != port {
			continue
		}
		if r.Proto != "" && r.Proto != proto {
			continue
		}
		return r.Action
	}
	return p.default_
}

// Summary returns a human-readable description of the policy rules.
func (p *Policy) Summary() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := fmt.Sprintf("default=%s rules=%d", p.default_, len(p.rules))
	return out
}
