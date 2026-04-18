// Package statestore persists the last known set of port states to disk
// so that portwatch can detect changes across restarts.
package statestore

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Store holds the most-recently observed port states.
type Store struct {
	mu      sync.RWMutex
	path    string
	states  []scanner.PortState
	updated time.Time
}

type persistedStore struct {
	Updated time.Time           `json:"updated"`
	States  []scanner.PortState `json:"states"`
}

// New returns an empty Store backed by path.
func New(path string) *Store {
	return &Store{path: path}
}

// Load reads previously persisted states from disk.
// Returns nil if the file does not exist.
func Load(path string) (*Store, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return New(path), nil
	}
	if err != nil {
		return nil, err
	}
	var p persistedStore
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &Store{path: path, states: p.States, updated: p.Updated}, nil
}

// Set replaces the current states and flushes to disk.
func (s *Store) Set(states []scanner.PortState) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states = states
	s.updated = time.Now()
	return s.flush()
}

// Get returns a copy of the current states.
func (s *Store) Get() []scanner.PortState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]scanner.PortState, len(s.states))
	copy(out, s.states)
	return out
}

// UpdatedAt returns when the store was last written.
func (s *Store) UpdatedAt() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.updated
}

func (s *Store) flush() error {
	p := persistedStore{Updated: s.updated, States: s.states}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o600)
}
