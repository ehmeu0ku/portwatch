// Package portdrift detects when a port's listening address drifts
// between scan cycles — for example, a service that was bound to
// 127.0.0.1 but is now bound to 0.0.0.0.
package portdrift

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// DriftEvent describes an address change for a given port.
type DriftEvent struct {
	Port    uint16
	Proto   string
	OldAddr string
	NewAddr string
}

func (d DriftEvent) String() string {
	return fmt.Sprintf("port %d/%s drifted %s -> %s", d.Port, d.Proto, d.OldAddr, d.NewAddr)
}

// Detector tracks the bound address of each port and reports when it changes.
type Detector struct {
	mu   sync.Mutex
	last map[string]string // key -> addr
}

// New returns an empty Detector.
func New() *Detector {
	return &Detector{last: make(map[string]string)}
}

func key(s scanner.PortState) string {
	return fmt.Sprintf("%d/%s", s.Port, s.Proto)
}

// Observe records the current address for each state and returns any
// DriftEvents for ports whose bound address has changed since the last call.
func (d *Detector) Observe(states []scanner.PortState) []DriftEvent {
	d.mu.Lock()
	defer d.mu.Unlock()

	var events []DriftEvent
	for _, s := range states {
		k := key(s)
		prev, seen := d.last[k]
		if seen && prev != s.Addr {
			events = append(events, DriftEvent{
				Port:    s.Port,
				Proto:   s.Proto,
				OldAddr: prev,
				NewAddr: s.Addr,
			})
		}
		d.last[k] = s.Addr
	}
	return events
}

// Forget removes tracking state for the given port/proto key so that the
// next Observe call treats it as a first-seen entry.
func (d *Detector) Forget(port uint16, proto string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.last, fmt.Sprintf("%d/%s", port, proto))
}

// Len returns the number of ports currently being tracked.
func (d *Detector) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.last)
}
