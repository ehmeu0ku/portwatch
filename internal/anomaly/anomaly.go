// Package anomaly detects unusual port activity based on historical scoring.
package anomaly

import (
	"fmt"
	"time"

	"github.com/user/portwatch/internal/scorecard"
)

// Level represents the anomaly severity level.
type Level int

const (
	LevelNone Level = iota
	LevelSuspicious
	LevelAnomalous
)

func (l Level) String() string {
	switch l {
	case LevelSuspicious:
		return "suspicious"
	case LevelAnomalous:
		return "anomalous"
	default:
		return "none"
	}
}

// Finding holds the result of an anomaly check.
type Finding struct {
	Key       string
	Level     Level
	Score     float64
	Reason    string
	DetectedAt time.Time
}

func (f Finding) String() string {
	return fmt.Sprintf("[%s] key=%s score=%.2f reason=%s", f.Level, f.Key, f.Score, f.Reason)
}

// Detector checks port activity against a scorecard.
type Detector struct {
	sc                *scorecard.Scorecard
	suspiciousThresh  float64
	anomalousThresh   float64
}

// New creates a Detector with the given scorecard and thresholds.
func New(sc *scorecard.Scorecard, suspicious, anomalous float64) *Detector {
	return &Detector{sc: sc, suspiciousThresh: suspicious, anomalousThresh: anomalous}
}

// Check evaluates the anomaly level for the given key.
func (d *Detector) Check(key string) (Finding, bool) {
	snap, ok := d.sc.Get(key)
	if !ok || snap.Seen == 0 {
		return Finding{}, false
	}
	score := snap.Score()
	f := Finding{Key: key, Score: score, DetectedAt: time.Now()}
	switch {
	case score >= d.anomalousThresh:
		f.Level = LevelAnomalous
		f.Reason = fmt.Sprintf("alert ratio %.0f%% exceeds anomalous threshold %.0f%%", score*100, d.anomalousThresh*100)
	case score >= d.suspiciousThresh:
		f.Level = LevelSuspicious
		f.Reason = fmt.Sprintf("alert ratio %.0f%% exceeds suspicious threshold %.0f%%", score*100, d.suspiciousThresh*100)
	default:
		return Finding{}, false
	}
	return f, true
}
