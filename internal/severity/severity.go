// Package severity classifies port events into alert levels based on
// configurable rules such as privileged port ranges and known service tags.
package severity

import "github.com/user/portwatch/internal/tagger"

// Level represents the urgency of a port event.
type Level int

const (
	Info    Level = iota // expected or low-interest activity
	Warning              // unusual but not necessarily malicious
	Critical             // high-priority: privileged port or untagged listener
)

// String returns a human-readable label for the level.
func (l Level) String() string {
	switch l {
	case Warning:
		return "WARNING"
	case Critical:
		return "CRITICAL"
	default:
		return "INFO"
	}
}

// Classifier assigns severity levels to port events.
type Classifier struct {
	privilegedMax uint16
	tagger        *tagger.Tagger
}

// New returns a Classifier. privilegedMax is the highest port number
// considered privileged (typically 1023). A nil tagger disables tag-based
// classification.
func New(privilegedMax uint16, t *tagger.Tagger) *Classifier {
	return &Classifier{privilegedMax: privilegedMax, tagger: t}
}

// Classify returns a Level for a newly-seen port number.
// Privileged ports are always Critical. Untagged high ports are Warning.
// Tagged high ports are Info.
func (c *Classifier) Classify(port uint16) Level {
	if port <= c.privilegedMax {
		return Critical
	}
	if c.tagger != nil {
		tag := c.tagger.Tag(port)
		if tag != "" {
			return Info
		}
	}
	return Warning
}
