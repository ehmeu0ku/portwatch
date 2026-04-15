// Package audit provides a persistent audit log of port events,
// recording new and gone port detections with timestamps for later review.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// EventKind describes the type of audit event.
type EventKind string

const (
	EventNew  EventKind = "NEW"
	EventGone EventKind = "GONE"
)

// Entry is a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Kind      EventKind `json:"kind"`
	Proto     string    `json:"proto"`
	Port      uint16    `json:"port"`
	Addr      string    `json:"addr"`
	PID       int       `json:"pid,omitempty"`
}

// Log is a thread-safe append-only audit log backed by a JSONL file.
type Log struct {
	mu   sync.Mutex
	path string
	f    *os.File
}

// Open opens (or creates) the audit log at the given path.
func Open(path string) (*Log, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o640)
	if err != nil {
		return nil, fmt.Errorf("audit: open %s: %w", path, err)
	}
	return &Log{path: path, f: f}, nil
}

// Record appends an entry to the audit log.
func (l *Log) Record(e Entry) error {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	_, err = fmt.Fprintf(l.f, "%s\n", data)
	return err
}

// Close closes the underlying file.
func (l *Log) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.f.Close()
}

// ReadAll reads all entries from the audit log file at path.
func ReadAll(path string) ([]Entry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("audit: read %s: %w", path, err)
	}
	var entries []Entry
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("audit: parse line: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
