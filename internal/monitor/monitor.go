// Package monitor provides the core polling loop that compares port states
// over time and dispatches alerts when unexpected changes are detected.
package monitor

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

// Config holds tunable parameters for the Monitor.
type Config struct {
	// Interval is how often the scanner is polled.
	Interval time.Duration
	// Baseline is the set of ports considered "expected". Ports outside this
	// set will trigger an alert.
	Baseline []scanner.PortState
}

// Monitor polls a Scanner on a fixed interval and notifies an Alerter when
// the set of listening ports changes relative to the last snapshot.
type Monitor struct {
	cfg     Config
	sc      scanner.Scanner
	al      *alert.Alerter
	last    []scanner.PortState
}

// New creates a Monitor with the provided dependencies.
func New(cfg Config, sc scanner.Scanner, al *alert.Alerter) *Monitor {
	return &Monitor{cfg: cfg, sc: sc, al: al}
}

// Run starts the polling loop and blocks until ctx is cancelled.
func (m *Monitor) Run(ctx context.Context) error {
	ticker := time.NewTicker(m.cfg.Interval)
	defer ticker.Stop()

	// Perform an initial scan so we have a baseline snapshot.
	if err := m.poll(); err != nil {
		log.Printf("portwatch: initial scan error: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := m.poll(); err != nil {
				log.Printf("portwatch: scan error: %v", err)
			}
		}
	}
}

// poll runs one scan cycle, diffs against the previous snapshot, and fires
// alerts for any new or disappeared ports.
func (m *Monitor) poll() error {
	current, err := m.sc.Scan()
	if err != nil {
		return err
	}
	current = scanner.DeduplicateStates(current)

	newPorts, gonePorts := diff(m.last, current)
	for _, p := range newPorts {
		m.al.Notify(p)
	}
	for _, p := range gonePorts {
		m.al.NotifyGone(p)
	}

	m.last = current
	return nil
}
