// Package snapshot provides functionality for capturing and persisting
// point-in-time views of active port states to disk.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot represents a point-in-time capture of port states.
type Snapshot struct {
	CapturedAt time.Time            `json:"captured_at"`
	States     []scanner.PortState  `json:"states"`
}

// New creates a new Snapshot from the given port states.
func New(states []scanner.PortState) *Snapshot {
	return &Snapshot{
		CapturedAt: time.Now().UTC(),
		States:     states,
	}
}

// Save writes the snapshot as JSON to the given file path.
func (s *Snapshot) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create file %q: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open file %q: %w", path, err)
	}
	defer f.Close()

	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &s, nil
}

// Summary returns a human-readable one-line description of the snapshot.
func (s *Snapshot) Summary() string {
	return fmt.Sprintf("snapshot captured at %s with %d port(s)",
		s.CapturedAt.Format(time.RFC3339), len(s.States))
}
