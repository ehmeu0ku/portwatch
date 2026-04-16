// Package formatter provides human-readable and structured formatting
// for port events emitted by the portwatch pipeline.
package formatter

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/portwatch/internal/correlator"
	"github.com/user/portwatch/internal/severity"
)

// Style controls the output format produced by Format.
type Style int

const (
	// StyleText produces a single human-readable line.
	StyleText Style = iota
	// StyleJSON produces a JSON object string.
	StyleJSON
)

// Formatter converts a correlator.Event into a printable string.
type Formatter struct {
	style Style
}

// New returns a Formatter that uses the given Style.
func New(s Style) *Formatter {
	return &Formatter{style: s}
}

// Format renders the event according to the configured style.
func (f *Formatter) Format(e correlator.Event) string {
	ts := e.Timestamp
	if ts.IsZero() {
		ts = time.Now().UTC()
	}

	switch f.style {
	case StyleJSON:
		return f.formatJSON(e, ts)
	default:
		return f.formatText(e, ts)
	}
}

func (f *Formatter) formatText(e correlator.Event, ts time.Time) string {
	sev := severityLabel(e.Severity)
	proto := strings.ToUpper(e.State.Proto)
	addr := fmt.Sprintf("%s:%d", e.State.IP, e.State.Port)
	process := "-"
	if e.State.Process != nil {
		process = e.State.Process.Name
	}
	return fmt.Sprintf("%s [%s] %s %s/%s pid=%d proc=%s tag=%s",
		ts.Format(time.RFC3339),
		sev,
		e.Kind,
		proto,
		addr,
		e.State.PID,
		process,
		e.Tag,
	)
}

func (f *Formatter) formatJSON(e correlator.Event, ts time.Time) string {
	process := ""
	if e.State.Process != nil {
		process = e.State.Process.Name
	}
	return fmt.Sprintf(
		`{"ts":%q,"kind":%q,"severity":%q,"proto":%q,"ip":%q,"port":%d,"pid":%d,"process":%q,"tag":%q}`,
		ts.Format(time.RFC3339),
		e.Kind,
		severityLabel(e.Severity),
		strings.ToUpper(e.State.Proto),
		e.State.IP,
		e.State.Port,
		e.State.PID,
		process,
		e.Tag,
	)
}

func severityLabel(s severity.Level) string {
	switch s {
	case severity.Critical:
		return "CRITICAL"
	case severity.Warning:
		return "WARNING"
	default:
		return "INFO"
	}
}
