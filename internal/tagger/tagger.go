// Package tagger assigns human-readable tags to port states based on
// well-known port numbers and configurable rules.
package tagger

import "github.com/user/portwatch/internal/scanner"

// WellKnown maps common port numbers to service names.
var WellKnown = map[uint16]string{
	22:   "ssh",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	443:  "https",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}

// Tagger annotates port states with service tags.
type Tagger struct {
	custom map[uint16]string
}

// New returns a Tagger that merges well-known tags with any caller-supplied
// custom mappings. Custom entries take precedence over well-known ones.
func New(custom map[uint16]string) *Tagger {
	merged := make(map[uint16]string, len(WellKnown)+len(custom))
	for k, v := range WellKnown {
		merged[k] = v
	}
	for k, v := range custom {
		merged[k] = v
	}
	return &Tagger{custom: merged}
}

// Tag returns the service name for the given port state, or an empty string
// when no mapping is known.
func (t *Tagger) Tag(s scanner.PortState) string {
	if name, ok := t.custom[s.Port]; ok {
		return name
	}
	return ""
}

// TagAll annotates a slice of PortStates, returning a map of port → tag for
// every state that has a known tag.
func (t *Tagger) TagAll(states []scanner.PortState) map[uint16]string {
	out := make(map[uint16]string)
	for _, s := range states {
		if tag := t.Tag(s); tag != "" {
			out[s.Port] = tag
		}
	}
	return out
}
