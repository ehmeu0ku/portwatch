// Package labelstore maintains a persistent map of port keys to
// user-defined labels, allowing operators to annotate known listeners.
package labelstore

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Label holds metadata attached to a port.
type Label struct {
	Name    string `json:"name"`
	Comment string `json:"comment,omitempty"`
}

// Key uniquely identifies a port listener.
type Key struct {
	Proto string
	Port  uint16
}

func (k Key) String() string {
	return fmt.Sprintf("%s:%d", k.Proto, k.Port)
}

// Store is a thread-safe label registry.
type Store struct {
	mu     sync.RWMutex
	path   string
	labels map[string]Label
}

// New returns an empty Store backed by path.
func New(path string) *Store {
	return &Store{path: path, labels: make(map[string]Label)}
}

// Set assigns a label to the given key.
func (s *Store) Set(k Key, l Label) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.labels[k.String()] = l
}

// Get retrieves a label. ok is false when no label exists.
func (s *Store) Get(k Key) (Label, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	l, ok := s.labels[k.String()]
	return l, ok
}

// Delete removes a label.
func (s *Store) Delete(k Key) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.labels, k.String())
}

// Save persists the store to disk as JSON.
func (s *Store) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, err := json.MarshalIndent(s.labels, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o644)
}

// Load reads labels from disk, replacing any in-memory state.
func (s *Store) Load() error {
	b, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	tmp := make(map[string]Label)
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.labels = tmp
	return nil
}
