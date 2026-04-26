// Package portmap maintains a live mapping of port keys to their most
// recently observed scanner states. It provides a thread-safe registry that
// the monitor can query to answer "what is currently listening on port X?"
// without re-scanning.
package portmap

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Key uniquely identifies a listening endpoint.
type Key struct {
	Port  uint16
	Proto string
}

// String returns a human-readable representation of the key.
func (k Key) String() string {
	return fmt.Sprintf("%s/%d", k.Proto, k.Port)
}

// KeyFromState derives a Key from a scanner.PortState.
func KeyFromState(s scanner.PortState) Key {
	return Key{Port: s.Port, Proto: s.Proto}
}

// Map is a thread-safe registry of currently observed port states.
type Map struct {
	mu     sync.RWMutex
	entries map[Key]scanner.PortState
}

// New returns an empty Map.
func New() *Map {
	return &Map{
		entries: make(map[Key]scanner.PortState),
	}
}

// Set stores or replaces the state for the given key.
func (m *Map) Set(s scanner.PortState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[KeyFromState(s)] = s
}

// Delete removes the entry for the given key, if present.
func (m *Map) Delete(k Key) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, k)
}

// Get returns the state for a key and whether it was found.
func (m *Map) Get(k Key) (scanner.PortState, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.entries[k]
	return s, ok
}

// Len returns the number of currently tracked entries.
func (m *Map) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.entries)
}

// Snapshot returns a shallow copy of all current entries.
func (m *Map) Snapshot() []scanner.PortState {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]scanner.PortState, 0, len(m.entries))
	for _, s := range m.entries {
		out = append(out, s)
	}
	return out
}

// Apply replaces the entire map contents with the provided slice, removing
// any keys that are no longer present. It returns the sets of added and
// removed keys so callers can react to changes.
func (m *Map) Apply(states []scanner.PortState) (added, removed []Key) {
	m.mu.Lock()
	defer m.mu.Unlock()

	next := make(map[Key]scanner.PortState, len(states))
	for _, s := range states {
		next[KeyFromState(s)] = s
	}

	for k := range next {
		if _, exists := m.entries[k]; !exists {
			added = append(added, k)
		}
	}
	for k := range m.entries {
		if _, exists := next[k]; !exists {
			removed = append(removed, k)
		}
	}

	m.entries = next
	return added, removed
}
