// Package reaper periodically removes stale snapshot and audit files
// that exceed a configured age, keeping disk usage bounded.
package reaper

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Reaper removes files older than MaxAge from Dir on every Interval tick.
type Reaper struct {
	Dir      string
	MaxAge   time.Duration
	Interval time.Duration
	Glob     string
	log      *log.Logger
}

// New returns a Reaper ready to run.
func New(dir string, maxAge, interval time.Duration, glob string, logger *log.Logger) *Reaper {
	if logger == nil {
		logger = log.New(os.Stderr, "[reaper] ", log.LstdFlags)
	}
	return &Reaper{
		Dir:      dir,
		MaxAge:   maxAge,
		Interval: interval,
		Glob:     glob,
		log:      logger,
	}
}

// Run blocks until ctx is cancelled, pruning files on each tick.
func (r *Reaper) Run(ctx context.Context) {
	ticker := time.NewTicker(r.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.Prune()
		}
	}
}

//etes files matching Glob inside Dir that are older than MaxAge.
// It returns the number of files removed.
func (r *Reaper) Prune() int {
	pattern := filepath.Join(r.Dir, r.Glob)
	mat(pattern)
	if err != nil {
		r.log.Printf("glob error: %v", err)
		return 0
	}
	cutoff := time.Now().Add(-r.MaxAge)
	removed := 0
	for _, path := range matches {
		info, err := os.Stat(path)
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			if err := os.Remove(path); err != nil {
				r.log.Printf("remove %s: %v", path, err)
				continue
			}
			r.log.Printf("removed stale file %s", path)
			removed++
		}
	}
	return removed
}
