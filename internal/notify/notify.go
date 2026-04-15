// Package notify provides pluggable notification channels for portwatch alerts.
package notify

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Message holds the data for a single notification event.
type Message struct {
	Level     Level
	Title     string
	Body      string
	Timestamp time.Time
}

// Notifier is the interface implemented by all notification backends.
type Notifier interface {
	Send(msg Message) error
	Name() string
}

// LogNotifier writes notifications as structured lines to an io.Writer.
type LogNotifier struct {
	w io.Writer
}

// NewLogNotifier returns a LogNotifier that writes to w.
// If w is nil, os.Stdout is used.
func NewLogNotifier(w io.Writer) *LogNotifier {
	if w == nil {
		w = os.Stdout
	}
	return &LogNotifier{w: w}
}

// Name returns the backend identifier.
func (l *LogNotifier) Name() string { return "log" }

// Send formats msg and writes it to the underlying writer.
func (l *LogNotifier) Send(msg Message) error {
	_, err := fmt.Fprintf(
		l.w,
		"%s [%s] %s: %s\n",
		msg.Timestamp.UTC().Format(time.RFC3339),
		msg.Level,
		msg.Title,
		msg.Body,
	)
	return err
}
