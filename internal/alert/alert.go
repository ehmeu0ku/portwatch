// Package alert provides alerting functionality for unexpected port listeners.
package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Alert represents a single alerting event for a port state change.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Message   string
	State     scanner.PortState
}

// Alerter handles formatting and dispatching alerts.
type Alerter struct {
	out    io.Writer
	prefix string
}

// NewAlerter creates a new Alerter writing to the given writer.
// If w is nil, os.Stdout is used.
func NewAlerter(w io.Writer) *Alerter {
	if w == nil {
		w = os.Stdout
	}
	return &Alerter{out: w, prefix: "portwatch"}
}

// Notify dispatches an alert for a newly detected port listener.
func (a *Alerter) Notify(level Level, state scanner.PortState) Alert {
	al := Alert{
		Timestamp: time.Now(),
		Level:     level,
		Message:   fmt.Sprintf("unexpected listener detected on %s", state),
		State:     state,
	}
	a.write(al)
	return al
}

// NotifyGone dispatches an alert when a previously seen port listener disappears.
func (a *Alerter) NotifyGone(state scanner.PortState) Alert {
	al := Alert{
		Timestamp: time.Now(),
		Level:     LevelInfo,
		Message:   fmt.Sprintf("listener closed on %s", state),
		State:     state,
	}
	a.write(al)
	return al
}

func (a *Alerter) write(al Alert) {
	fmt.Fprintf(a.out, "[%s] %s %s\n",
		al.Timestamp.Format(time.RFC3339),
		al.Level,
		al.Message,
	)
}
