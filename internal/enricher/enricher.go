// Package enricher attaches process information to port states.
// It bridges the scanner and process packages so that alert messages
// can include the owning process name and PID alongside the port details.
package enricher

import (
	"fmt"

	"github.com/user/portwatch/internal/process"
	"github.com/user/portwatch/internal/scanner"
)

// EnrichedState wraps a scanner.PortState with optional process metadata.
type EnrichedState struct {
	scanner.PortState
	Process *process.Info
}

// String returns a human-readable representation that includes process info
// when available.
func (e EnrichedState) String() string {
	base := e.PortState.String()
	if e.Process == nil {
		return base
	}
	return fmt.Sprintf("%s [%s]", base, e.Process.String())
}

// Enricher resolves process information for port states.
type Enricher struct {
	resolver *process.Resolver
}

// New creates an Enricher backed by the given Resolver.
func New(r *process.Resolver) *Enricher {
	return &Enricher{resolver: r}
}

// Enrich looks up the process owning each PortState's inode and returns a
// slice of EnrichedStates. States whose inodes cannot be resolved still
// appear in the result; their Process field will be nil.
func (e *Enricher) Enrich(states []scanner.PortState) []EnrichedState {
	out := make([]EnrichedState, 0, len(states))
	for _, s := range states {
		es := EnrichedState{PortState: s}
		if info, err := e.resolver.LookupInode(s.Inode); err == nil {
			es.Process = info
		}
		out = append(out, es)
	}
	return out
}

// EnrichOne is a convenience wrapper for a single PortState.
func (e *Enricher) EnrichOne(s scanner.PortState) EnrichedState {
	es := EnrichedState{PortState: s}
	if info, err := e.resolver.LookupInode(s.Inode); err == nil {
		es.Process = info
	}
	return es
}
