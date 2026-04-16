// Package deadman provides a dead-man's switch that fires an alert when
// the monitor has not produced any scan results within a configurable window.
package deadman

import (
	"context"
	"sync"
	"time"
)

// AlertFunc is called when the dead-man's switch trips.
type AlertFunc func(missed time.Duration)

// Switch fires AlertFunc if Reset is not called within Window.
type Switch struct {
	mu      sync.Mutex
	window  time.Duration
	alert   AlertFunc
	last    time.Time
	timer   *time.Timer
	tripped bool
}

// New creates a Switch with the given window and alert callback.
func New(window time.Duration, fn AlertFunc) *Switch {
	return &Switch{
		window: window,
		alert:  fn,
		last:   time.Now(),
	}
}

// Start begins monitoring in the background; it stops when ctx is cancelled.
func (s *Switch) Start(ctx context.Context) {
	s.mu.Lock()
	s.timer = time.AfterFunc(s.window, func() { s.trip() })
	s.mu.Unlock()

	go func() {
		<-ctx.Done()
		s.mu.Lock()
		if s.timer != nil {
			s.timer.Stop()
		}
		s.mu.Unlock()
	}()
}

// Reset signals that the monitor is alive; it resets the countdown.
func (s *Switch) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.last = time.Now()
	s.tripped = false
	if s.timer != nil {
		s.timer.Reset(s.window)
	}
}

// Tripped reports whether the switch has fired.
func (s *Switch) Tripped() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tripped
}

func (s *Switch) trip() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.tripped {
		return
	}
	s.tripped = true
	missed := time.Since(s.last)
	go s.alert(missed)
}
