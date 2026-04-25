// Package portmatch evaluates whether a port state matches a set of
// user-defined match expressions. Expressions support exact port numbers,
// port ranges, protocol filters, and tag predicates.
package portmatch

import (
	"fmt"
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

// Expr is a single match expression.
type Expr struct {
	// Port is an exact port number to match (0 = any).
	Port uint16
	// PortMin/PortMax define an inclusive range (both 0 = any).
	PortMin uint16
	PortMax uint16
	// Proto restricts the match to "tcp" or "udp" (empty = any).
	Proto string
	// Tag requires the port state to carry this tag (empty = any).
	Tag string
}

// Match returns true when state satisfies all non-zero fields in e.
func (e Expr) Match(s scanner.PortState) bool {
	if e.Proto != "" && !strings.EqualFold(s.Proto, e.Proto) {
		return false
	}
	if e.Tag != "" {
		found := false
		for _, t := range s.Tags {
			if strings.EqualFold(t, e.Tag) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if e.Port != 0 {
		return s.Port == e.Port
	}
	if e.PortMin != 0 || e.PortMax != 0 {
		return s.Port >= e.PortMin && s.Port <= e.PortMax
	}
	return true
}

// Matcher holds a list of expressions and tests states against them.
type Matcher struct {
	exprs []Expr
}

// New creates a Matcher from the provided expressions.
func New(exprs []Expr) *Matcher {
	return &Matcher{exprs: exprs}
}

// AnyMatch returns true if at least one expression matches the state.
func (m *Matcher) AnyMatch(s scanner.PortState) bool {
	for _, e := range m.exprs {
		if e.Match(s) {
			return true
		}
	}
	return false
}

// AllMatch returns true if every expression matches the state.
func (m *Matcher) AllMatch(s scanner.PortState) bool {
	for _, e := range m.exprs {
		if !e.Match(s) {
			return false
		}
	}
	return true
}

// String returns a human-readable description of the matcher.
func (m *Matcher) String() string {
	parts := make([]string, 0, len(m.exprs))
	for _, e := range m.exprs {
		parts = append(parts, exprString(e))
	}
	return fmt.Sprintf("Matcher[%s]", strings.Join(parts, ", "))
}

func exprString(e Expr) string {
	var b strings.Builder
	if e.Proto != "" {
		b.WriteString(e.Proto + ":")
	}
	switch {
	case e.Port != 0:
		fmt.Fprintf(&b, "%d", e.Port)
	case e.PortMin != 0 || e.PortMax != 0:
		fmt.Fprintf(&b, "%d-%d", e.PortMin, e.PortMax)
	default:
		b.WriteString("*")
	}
	if e.Tag != "" {
		b.WriteString("#" + e.Tag)
	}
	return b.String()
}
