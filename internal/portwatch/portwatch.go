// Package portwatch provides a port-to-process ownership tracker that
// records which process first claimed a given port and flags ownership changes.
package portwatch

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Entry records the first-seen process information for a port.
type Entry struct {
	Port    uint16
	Proto   string
	PID     int
	Comm    string
}

// ChangeKind describes how ownership changed.
type ChangeKind string

const (
	OwnershipClaimed  ChangeKind = "claimed"
	OwnershipChanged  ChangeKind = "changed"
	OwnershipReleased ChangeKind = "released"
)

// Change represents an ownership event for a port.
type Change struct {
	Kind    ChangeKind
	Prev    *Entry
	Current *Entry
}

func (c Change) String() string {
	switch c.Kind {
	case OwnershipClaimed:
		return fmt.Sprintf("claimed %s/%d by pid %d (%s)", c.Current.Proto, c.Current.Port, c.Current.PID, c.Current.Comm)
	case OwnershipChanged:
		return fmt.Sprintf("changed %s/%d: pid %d (%s) -> pid %d (%s)", c.Current.Proto, c.Current.Port, c.Prev.PID, c.Prev.Comm, c.Current.PID, c.Current.Comm)
	case OwnershipReleased:
		return fmt.Sprintf("released %s/%d (was pid %d %s)", c.Prev.Proto, c.Prev.Port, c.Prev.PID, c.Prev.Comm)
	default:
		return "unknown change"
	}
}

type key struct {
	Port  uint16
	Proto string
}

// Tracker maintains ownership state across scans.
type Tracker struct {
	mu      sync.Mutex
	owners  map[key]*Entry
}

// New creates an empty Tracker.
func New() *Tracker {
	return &Tracker{owners: make(map[key]*Entry)}
}

// Update compares current port states against recorded owners and returns any changes.
func (t *Tracker) Update(states []scanner.PortState) []Change {
	t.mu.Lock()
	defer t.mu.Unlock()

	seen := make(map[key]bool)
	var changes []Change

	for _, s := range states {
		if s.Process == nil {
			continue
		}
		k := key{Port: s.Port, Proto: s.Proto}
		seen[k] = true
		current := &Entry{Port: s.Port, Proto: s.Proto, PID: s.Process.PID, Comm: s.Process.Comm}
		prev, exists := t.owners[k]
		if !exists {
			changes = append(changes, Change{Kind: OwnershipClaimed, Current: current})
		} else if prev.PID != current.PID {
			changes = append(changes, Change{Kind: OwnershipChanged, Prev: prev, Current: current})
		}
		t.owners[k] = current
	}

	for k, entry := range t.owners {
		if !seen[k] {
			changes = append(changes, Change{Kind: OwnershipReleased, Prev: entry})
			delete(t.owners, k)
		}
	}
	return changes
}
