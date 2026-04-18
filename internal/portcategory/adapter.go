package portcategory

import (
	"github.com/yourusername/portwatch/internal/correlator"
)

// Annotator wraps a Classifier and attaches category labels to correlator
// events as a tag string, e.g. "category:web".
type Annotator struct {
	cl *Classifier
}

// NewAnnotator returns an Annotator backed by the given Classifier.
func NewAnnotator(cl *Classifier) *Annotator {
	return &Annotator{cl: cl}
}

// Annotate returns a copy of the event with a category tag appended.
func (a *Annotator) Annotate(ev correlator.Event) correlator.Event {
	cat := a.cl.Classify(uint16(ev.State.Port))
	if cat == Unknown {
		return ev
	}
	tag := "category:" + string(cat)
	for _, t := range ev.Tags {
		if t == tag {
			return ev
		}
	}
	copy := ev
	copy.Tags = append(append([]string(nil), ev.Tags...), tag)
	return copy
}
