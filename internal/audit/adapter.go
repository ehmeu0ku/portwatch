package audit

import (
	"github.com/your-org/portwatch/internal/scanner"
)

// FromPortState converts a scanner.PortState into an audit Entry.
// kind must be EventNew or EventGone.
func FromPortState(kind EventKind, s scanner.PortState) Entry {
	return Entry{
		Kind:  kind,
		Proto: s.Proto,
		Port:  s.Port,
		Addr:  s.Addr,
		PID:   s.PID,
	}
}

// Recorder wraps a Log and exposes convenience methods that accept
// scanner.PortState values directly.
type Recorder struct {
	log *Log
}

// NewRecorder creates a Recorder backed by log.
func NewRecorder(log *Log) *Recorder {
	return &Recorder{log: log}
}

// RecordNew appends a NEW event for the given state.
func (r *Recorder) RecordNew(s scanner.PortState) error {
	return r.log.Record(FromPortState(EventNew, s))
}

// RecordGone appends a GONE event for the given state.
func (r *Recorder) RecordGone(s scanner.PortState) error {
	return r.log.Record(FromPortState(EventGone, s))
}
