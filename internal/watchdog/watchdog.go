// Package watchdog provides a self-monitoring component that detects
// when the scanner or monitor loop has stalled and triggers recovery.
package watchdog

import (
	"context"
	"log"
	"sync"
	"time"
)

// Watchdog monitors a heartbeat channel and fires a callback when the
// heartbeat has not been received within the configured timeout.
type Watchdog struct {
	timeout   time.Duration
	onTimeout func()
	heartbeat chan struct{}
	mu        sync.Mutex
	running   bool
}

// New creates a new Watchdog with the given timeout duration and timeout
// callback. The callback is invoked in its own goroutine.
func New(timeout time.Duration, onTimeout func()) *Watchdog {
	return &Watchdog{
		timeout:   timeout,
		onTimeout: onTimeout,
		heartbeat:  make(chan struct{}, 1),
	}
}

// Beat signals the watchdog that the monitored component is alive.
// Calling Beat resets the internal timer.
func (w *Watchdog) Beat() {
	select {
	case w.heartbeat <- struct{}{}:
	default:
		// channel already has a pending beat; no need to queue another
	}
}

// Start begins the watchdog loop. It blocks until ctx is cancelled.
func (w *Watchdog) Start(ctx context.Context) {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return
	}
	w.running = true
	w.mu.Unlock()

	timer := time.NewTimer(w.timeout)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			w.mu.Lock()
			w.running = false
			w.mu.Unlock()
			return
		case <-w.heartbeat:
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(w.timeout)
		case <-timer.C:
			log.Printf("[watchdog] timeout after %s — triggering recovery", w.timeout)
			go w.onTimeout()
			timer.Reset(w.timeout)
		}
	}
}

// Running reports whether the watchdog loop is currently active.
func (w *Watchdog) Running() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.running
}
