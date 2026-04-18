package anomaly

import (
	"fmt"

	"github.com/user/portwatch/internal/correlator"
	"github.com/user/portwatch/internal/scorecard"
)

// KeyFromEvent derives a scorecard key from a correlator event.
func KeyFromEvent(e correlator.Event) string {
	return fmt.Sprintf("%s:%d", e.State.Proto, e.State.Port)
}

// Recorder wraps a Scorecard and records events from the correlator pipeline.
type Recorder struct {
	sc *scorecard.Scorecard
}

// NewRecorder creates a Recorder backed by sc.
func NewRecorder(sc *scorecard.Scorecard) *Recorder {
	return &Recorder{sc: sc}
}

// Record updates the scorecard for the event's key.
// An event is counted as an alert when its kind is "new".
func (r *Recorder) Record(e correlator.Event) {
	key := KeyFromEvent(e)
	isAlert := e.Kind == "new"
	r.sc.Record(key, isAlert)
}
