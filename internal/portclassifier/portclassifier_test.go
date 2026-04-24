package portclassifier_test

import (
	"testing"

	"github.com/user/portwatch/internal/portclassifier"
	"github.com/user/portwatch/internal/tagger"
)

func newClassifier(highPorts []uint16) *portclassifier.Classifier {
	t := tagger.New(nil) // default well-known mappings only
	return portclassifier.New(t, highPorts)
}

func TestClassificationString(t *testing.T) {
	cases := []struct {
		c    portclassifier.Classification
		want string
	}{
		{portclassifier.Benign, "benign"},
		{portclassifier.Elevated, "elevated"},
		{portclassifier.High, "high"},
	}
	for _, tc := range cases {
		if got := tc.c.String(); got != tc.want {
			t.Errorf("String() = %q, want %q", got, tc.want)
		}
	}
}

func TestTaggedPortIsBenign(t *testing.T) {
	c := newClassifier(nil)
	// port 80 is tagged as "http" by the default tagger
	if got := c.Classify(80, "tcp"); got != portclassifier.Benign {
		t.Errorf("port 80 tcp: got %s, want benign", got)
	}
}

func TestUntaggedPrivilegedPortIsElevated(t *testing.T) {
	c := newClassifier(nil)
	// port 999 is below 1024 and has no well-known tag
	if got := c.Classify(999, "tcp"); got != portclassifier.Elevated {
		t.Errorf("port 999 tcp: got %s, want elevated", got)
	}
}

func TestUntaggedHighPortIsBenign(t *testing.T) {
	c := newClassifier(nil)
	if got := c.Classify(49152, "tcp"); got != portclassifier.Benign {
		t.Errorf("port 49152 tcp: got %s, want benign", got)
	}
}

func TestExplicitHighPortIsHigh(t *testing.T) {
	c := newClassifier([]uint16{4444})
	if got := c.Classify(4444, "tcp"); got != portclassifier.High {
		t.Errorf("port 4444 tcp: got %s, want high", got)
	}
}

func TestHighPortOverridesTag(t *testing.T) {
	// Even a well-known port should be High when explicitly registered.
	c := newClassifier([]uint16{80})
	if got := c.Classify(80, "tcp"); got != portclassifier.High {
		t.Errorf("port 80 tcp (forced high): got %s, want high", got)
	}
}

func TestMultipleHighPorts(t *testing.T) {
	c := newClassifier([]uint16{1337, 31337})
	for _, p := range []uint16{1337, 31337} {
		if got := c.Classify(p, "tcp"); got != portclassifier.High {
			t.Errorf("port %d: got %s, want high", p, got)
		}
	}
}
