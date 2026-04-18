// Package tunables provides runtime-adjustable parameters for portwatch.
// Values can be updated without restarting the daemon via a reload signal.
package tunables

import (
	"sync"
	"time"
)

// Tunables holds runtime-adjustable operational parameters.
type Tunables struct {
	mu             sync.RWMutex
	scanInterval   time.Duration
	alertCooldown  time.Duration
	maxHistorySize int
}

// Defaults returns a Tunables instance populated with sensible defaults.
func Defaults() *Tunables {
	return &Tunables{
		scanInterval:   5 * time.Second,
		alertCooldown:  30 * time.Second,
		maxHistorySize: 500,
	}
}

// ScanInterval returns the current scan interval.
func (t *Tunables) ScanInterval() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.scanInterval
}

// SetScanInterval updates the scan interval. Returns false if d < 1s.
func (t *Tunables) SetScanInterval(d time.Duration) bool {
	if d < time.Second {
		return false
	}
	t.mu.Lock()
	t.scanInterval = d
	t.mu.Unlock()
	return true
}

// AlertCooldown returns the current alert cooldown duration.
func (t *Tunables) AlertCooldown() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.alertCooldown
}

// SetAlertCooldown updates the alert cooldown. Returns false if d is negative.
func (t *Tunables) SetAlertCooldown(d time.Duration) bool {
	if d < 0 {
		return false
	}
	t.mu.Lock()
	t.alertCooldown = d
	t.mu.Unlock()
	return true
}

// MaxHistorySize returns the maximum number of history entries to retain.
func (t *Tunables) MaxHistorySize() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.maxHistorySize
}

// SetMaxHistorySize updates the history size cap. Returns false if n < 1.
func (t *Tunables) SetMaxHistorySize(n int) bool {
	if n < 1 {
		return false
	}
	t.mu.Lock()
	t.maxHistorySize = n
	t.mu.Unlock()
	return true
}
