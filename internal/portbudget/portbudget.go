// Package portbudget enforces a maximum number of concurrently observed open
// ports per protocol. When the tracked count exceeds the configured ceiling the
// Budget.Exceeds method returns true, allowing callers to raise an alert or
// apply back-pressure.
package portbudget

import (
	"fmt"
	"sync"
)

// Budget tracks the number of open ports per protocol and reports when a
// configured ceiling is breached.
type Budget struct {
	mu      sync.Mutex
	counts  map[string]int
	ceiling int
}

// New returns a Budget that will flag any protocol whose open-port count
// exceeds ceiling. A ceiling of zero disables enforcement (Exceeds always
// returns false).
func New(ceiling int) *Budget {
	return &Budget{
		counts:  make(map[string]int),
		ceiling: ceiling,
	}
}

// Observe records proto as having one additional open port. It returns true if
// the updated count now exceeds the budget ceiling.
func (b *Budget) Observe(proto string) bool {
	if b.ceiling <= 0 {
		return false
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.counts[proto]++
	return b.counts[proto] > b.ceiling
}

// Release decrements the open-port count for proto. It is a no-op if the
// count is already zero.
func (b *Budget) Release(proto string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.counts[proto] > 0 {
		b.counts[proto]--
	}
}

// Count returns the current open-port count for proto.
func (b *Budget) Count(proto string) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.counts[proto]
}

// Exceeds reports whether proto's current count is over the ceiling without
// modifying any state.
func (b *Budget) Exceeds(proto string) bool {
	if b.ceiling <= 0 {
		return false
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.counts[proto] > b.ceiling
}

// Reset clears all counts for proto.
func (b *Budget) Reset(proto string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.counts, proto)
}

// Summary returns a human-readable description of every protocol whose count
// exceeds the ceiling.
func (b *Budget) Summary() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.ceiling <= 0 {
		return "budget enforcement disabled"
	}
	out := ""
	for proto, n := range b.counts {
		if n > b.ceiling {
			out += fmt.Sprintf("%s:%d/%d ", proto, n, b.ceiling)
		}
	}
	if out == "" {
		return "within budget"
	}
	return out[:len(out)-1]
}
