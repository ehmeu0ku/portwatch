// Package baseline manages the known-good set of listening ports.
// It persists state to disk so portwatch can detect deviations across restarts.
package baseline

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry records a single approved port listener in the baseline.
type Entry struct {
	Proto   string    `json:"proto"`
	Address string    `json:"address"`
	Port    uint16    `json:"port"`
	AddedAt time.Time `json:"added_at"`
}

// Baseline holds the approved set of port states.
type Baseline struct {
	mu      sync.RWMutex
	entries map[string]Entry
	path    string
}

// New creates an empty Baseline backed by the given file path.
func New(path string) *Baseline {
	return &Baseline{
		entries: make(map[string]Entry),
		path:    path,
	}
}

// Load reads the baseline from disk. Returns nil if the file does not exist.
func Load(path string) (*Baseline, error) {
	b := New(path)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return b, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	for _, e := range entries {
		b.entries[entryKey(e.Proto, e.Address, e.Port)] = e
	}
	return b, nil
}

// Save persists the current baseline to disk.
func (b *Baseline) Save() error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	list := make([]Entry, 0, len(b.entries))
	for _, e := range b.entries {
		list = append(list, e)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(b.path, data, 0o600)
}

// Add approves a port state, adding it to the baseline.
func (b *Baseline) Add(s scanner.PortState) {
	b.mu.Lock()
	defer b.mu.Unlock()
	k := entryKey(s.Proto, s.Address, s.Port)
	b.entries[k] = Entry{Proto: s.Proto, Address: s.Address, Port: s.Port, AddedAt: time.Now()}
}

// Contains reports whether the given port state is in the baseline.
func (b *Baseline) Contains(s scanner.PortState) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, ok := b.entries[entryKey(s.Proto, s.Address, s.Port)]
	return ok
}

// Entries returns a snapshot of all baseline entries.
func (b *Baseline) Entries() []Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]Entry, 0, len(b.entries))
	for _, e := range b.entries {
		out = append(out, e)
	}
	return out
}

func entryKey(proto, addr string, port uint16) string {
	return proto + "|" + addr + "|" + string(rune(port))
}
