// Package correlator links port events to enriched process information
// and assigns severity, producing a unified PortEvent ready for dispatch.
package correlator

import (
	"time"

	"github.com/user/portwatch/internal/enricher"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/severity"
	"github.com/user/portwatch/internal/tagger"
)

// Kind describes whether a port appeared or disappeared.
type Kind string

const (
	KindNew  Kind = "NEW"
	KindGone Kind = "GONE"
)

// PortEvent is the fully-correlated record produced for each state change.
type PortEvent struct {
	Kind      Kind
	State     enricher.EnrichedState
	Severity  severity.Level
	Tag       string
	Timestamp time.Time
}

// Correlator enriches raw scanner states and assigns severity + tags.
type Correlator struct {
	enricher  *enricher.Enricher
	tagger    *tagger.Tagger
	severity  *severity.Scorer
}

// New constructs a Correlator from its dependencies.
func New(e *enricher.Enricher, t *tagger.Tagger, s *severity.Scorer) *Correlator {
	return &Correlator{enricher: e, tagger: t, severity: s}
}

// Correlate converts a raw PortState and a Kind into a PortEvent.
func (c *Correlator) Correlate(kind Kind, ps scanner.PortState) PortEvent {
	es := c.enricher.Enrich(ps)
	tag := c.tagger.Tag(ps)
	lvl := c.severity.Score(ps, tag)
	return PortEvent{
		Kind:      kind,
		State:     es,
		Severity:  lvl,
		Tag:       tag,
		Timestamp: time.Now(),
	}
}
