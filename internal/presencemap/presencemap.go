// Package presencemap tracks which ports have been continuously present
// across consecutive scans, enabling stable-state detection.
package presencemap

import "sync"

// PresenceMap records how many consecutive scans each port key has been seen.
type PresenceMap struct {
	mu      sync.Mutex
	counts  map[string]int
	threshold int
}

// New returns a PresenceMap that considers a port stable after threshold
// consecutive observations.
func New(threshold int) *PresenceMap {
	if threshold < 1 {
		threshold = 1
	}
	return &PresenceMap{
		counts:    make(map[string]int),
		threshold: threshold,
	}
}

// Observe increments the consecutive-seen counter for key and returns true
// once the counter reaches the configured threshold.
func (p *PresenceMap) Observe(key string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.counts[key]++
	return p.counts[key] >= p.threshold
}

// Forget removes key from the map, resetting its counter.
func (p *PresenceMap) Forget(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.counts, key)
}

// Count returns the current consecutive-seen count for key.
func (p *PresenceMap) Count(key string) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.counts[key]
}

// Stable returns true if key has already reached the threshold.
func (p *PresenceMap) Stable(key string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.counts[key] >= p.threshold
}

// Len returns the number of keys currently tracked.
func (p *PresenceMap) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.counts)
}
