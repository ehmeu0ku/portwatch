// Package fingerprint generates stable identifiers for port states,
// enabling deduplication and change detection across scan cycles.
package fingerprint

import (
	"fmt"
	"hash/fnv"

	"github.com/user/portwatch/internal/scanner"
)

// Fingerprint is a stable hash of a port state's key attributes.
type Fingerprint uint64

// String returns a hex representation of the fingerprint.
func (f Fingerprint) String() string {
	return fmt.Sprintf("%016x", uint64(f))
}

// Of computes a fingerprint for the given PortState using its
// protocol, address, and port number.
func Of(s scanner.PortState) Fingerprint {
	h := fnv.New64a()
	_, _ = fmt.Fprintf(h, "%s:%s:%d", s.Proto, s.IP, s.Port)
	return Fingerprint(h.Sum64())
}

// Map holds fingerprints keyed by their string representation,
// mapping back to the originating PortState.
type Map map[string]scanner.PortState

// Build constructs a Map from a slice of PortStates.
func Build(states []scanner.PortState) Map {
	m := make(Map, len(states))
	for _, s := range states {
		f := Of(s)
		m[f.String()] = s
	}
	return m
}

// Added returns states present in next but not in prev.
func Added(prev, next Map) []scanner.PortState {
	var out []scanner.PortState
	for k, s := range next {
		if _, ok := prev[k]; !ok {
			out = append(out, s)
		}
	}
	return out
}

// Removed returns states present in prev but not in next.
func Removed(prev, next Map) []scanner.PortState {
	var out []scanner.PortState
	for k, s := range prev {
		if _, ok := next[k]; !ok {
			out = append(out, s)
		}
	}
	return out
}
