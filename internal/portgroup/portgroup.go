// Package portgroup groups related ports under a named label for
// cleaner alerting and filtering.
package portgroup

import "sync"

// Group is a named collection of port numbers.
type Group struct {
	Name  string
	Ports map[uint16]struct{}
}

// Store holds named port groups and allows membership queries.
type Store struct {
	mu     sync.RWMutex
	groups map[string]*Group
}

// New returns an empty Store.
func New() *Store {
	return &Store{groups: make(map[string]*Group)}
}

// Define registers a named group with the given ports.
// Calling Define with an existing name replaces the group.
func (s *Store) Define(name string, ports []uint16) {
	g := &Group{Name: name, Ports: make(map[uint16]struct{}, len(ports))}
	for _, p := range ports {
		g.Ports[p] = struct{}{}
	}
	s.mu.Lock()
	s.groups[name] = g
	s.mu.Unlock()
}

// Lookup returns the names of all groups that contain the given port.
func (s *Store) Lookup(port uint16) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var names []string
	for name, g := range s.groups {
		if _, ok := g.Ports[port]; ok {
			names = append(names, name)
		}
	}
	return names
}

// Contains reports whether the named group includes the given port.
func (s *Store) Contains(name string, port uint16) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	g, ok := s.groups[name]
	if !ok {
		return false
	}
	_, found := g.Ports[port]
	return found
}

// Names returns all defined group names.
func (s *Store) Names() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, 0, len(s.groups))
	for n := range s.groups {
		out = append(out, n)
	}
	return out
}
