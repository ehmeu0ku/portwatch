// Package pagerule evaluates a set of user-defined rules against a
// correlator.Event and decides whether a desktop/webhook page should fire.
package pagerule

import (
	"strings"

	"github.com/user/portwatch/internal/correlator"
	"github.com/user/portwatch/internal/severity"
)

// Action describes what to do when a rule matches.
type Action string

const (
	ActionPage    Action = "page"
	ActionSuppress Action = "suppress"
)

// Rule is a single matching rule.
type Rule struct {
	// MinSeverity is the minimum severity level that triggers this rule.
	MinSeverity severity.Level
	// Ports, if non-empty, restricts the rule to specific ports.
	Ports []uint16
	// Tags, if non-empty, requires at least one tag to match.
	Tags   []string
	Action Action
}

// Evaluator checks an event against an ordered list of rules.
type Evaluator struct {
	rules []Rule
}

// New returns an Evaluator with the provided rules.
func New(rules []Rule) *Evaluator {
	return &Evaluator{rules: rules}
}

// Evaluate returns the Action of the first matching rule, or ActionSuppress
// if no rule matches.
func (e *Evaluator) Evaluate(ev correlator.Event) Action {
	for _, r := range e.rules {
		if ev.Severity < r.MinSeverity {
			continue
		}
		if len(r.Ports) > 0 && !containsPort(r.Ports, ev.State.Port) {
			continue
		}
		if len(r.Tags) > 0 && !hasTag(r.Tags, ev.Tags) {
			continue
		}
		return r.Action
	}
	return ActionSuppress
}

func containsPort(ports []uint16, port uint16) bool {
	for _, p := range ports {
		if p == port {
			return true
		}
	}
	return false
}

func hasTag(want []string, have []string) bool {
	for _, w := range want {
		for _, h := range have {
			if strings.EqualFold(w, h) {
				return true
			}
		}
	}
	return false
}
