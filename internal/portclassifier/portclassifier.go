// Package portclassifier assigns a risk classification to a port based on
// its number, protocol, and associated tags. Classifications range from
// Benign through Elevated to High, and are intended to feed downstream
// alerting and dispatch pipelines.
package portclassifier

import "github.com/user/portwatch/internal/tagger"

// Classification represents the risk level assigned to a port.
type Classification int

const (
	Benign   Classification = iota // expected, well-understood service
	Elevated                       // unusual but not necessarily malicious
	High                           // unexpected or high-risk listener
)

// String returns a human-readable label for the classification.
func (c Classification) String() string {
	switch c {
	case Benign:
		return "benign"
	case Elevated:
		return "elevated"
	case High:
		return "high"
	default:
		return "unknown"
	}
}

// Classifier assigns risk classifications to ports.
type Classifier struct {
	tagger    *tagger.Tagger
	highPorts map[uint16]struct{}
}

// New returns a Classifier backed by the provided Tagger.
// Additional ports can be registered as unconditionally High-risk via
// the optional highPorts slice.
func New(t *tagger.Tagger, highPorts []uint16) *Classifier {
	hp := make(map[uint16]struct{}, len(highPorts))
	for _, p := range highPorts {
		hp[p] = struct{}{}
	}
	return &Classifier{tagger: t, highPorts: hp}
}

// Classify returns the Classification for the given port number and protocol.
// Ports below 1024 are always at least Elevated. Ports explicitly registered
// as high-risk are always High. Tagged ports are considered Benign.
func (c *Classifier) Classify(port uint16, proto string) Classification {
	if _, ok := c.highPorts[port]; ok {
		return High
	}

	tag := c.tagger.Tag(port, proto)
	if tag != "" {
		return Benign
	}

	if port < 1024 {
		return Elevated
	}

	return Benign
}
