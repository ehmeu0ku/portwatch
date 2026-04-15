// Package filter provides port filtering utilities for portwatch.
// It allows callers to suppress known-safe ports from scan results
// before they are passed to the alerter.
package filter

import (
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
)

// Filter removes port states that should be ignored according to the
// provided configuration (ignored ports and loopback-only filtering).
type Filter struct {
	cfg *config.Config
}

// New returns a new Filter backed by cfg.
func New(cfg *config.Config) *Filter {
	return &Filter{cfg: cfg}
}

// Apply returns a new slice containing only the states that are NOT
// suppressed by the current configuration.
func (f *Filter) Apply(states []scanner.PortState) []scanner.PortState {
	out := make([]scanner.PortState, 0, len(states))
	for _, s := range states {
		if f.suppress(s) {
			continue
		}
		out = append(out, s)
	}
	return out
}

// suppress returns true when the given state should be hidden from alerts.
func (f *Filter) suppress(s scanner.PortState) bool {
	// Drop loopback-only listeners when the config asks for it.
	if f.cfg.IgnoreLoopback && scanner.IsLoopback(s.IP) {
		return true
	}
	// Drop explicitly ignored ports.
	if f.cfg.IsIgnored(s.Port) {
		return true
	}
	return false
}
