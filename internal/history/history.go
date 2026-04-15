// Package history tracks port state changes over time,
// allowing portwatch to report trends and recurring events.
package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry records a single port state change event.
type Entry struct {
	Timestamp time.Time          `json:"timestamp"`
	State     scanner.PortState  `json:"state"`
	Event     string             `json:"event"` // "new" or "gone"
}

// History holds a bounded log of port change events.
type History struct {
	mu      sync.RWMutex
	entries []Entry
	maxSize int
}

// New creates a History with the given maximum number of entries.
func New(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 500
	}
	return &History{maxSize: maxSize}
}

// Record appends a new event to the history, evicting the oldest if full.
func (h *History) Record(state scanner.PortState, event string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	e := Entry{Timestamp: time.Now(), State: state, Event: event}
	if len(h.entries) >= h.maxSize {
		h.entries = h.entries[1:]
	}
	h.entries = append(h.entries, e)
}

// Entries returns a snapshot of all recorded entries.
func (h *History) Entries() []Entry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]Entry, len(h.entries))
	copy(out, h.entries)
	return out
}

// Save persists the history to a JSON file at path.
func (h *History) Save(path string) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(h.entries)
}

// Load reads history entries from a JSON file, replacing any existing entries.
func (h *History) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	var entries []Entry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return err
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = entries
	return nil
}
