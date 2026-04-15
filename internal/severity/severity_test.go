package severity_test

import (
	"testing"

	"github.com/user/portwatch/internal/severity"
	"github.com/user/portwatch/internal/tagger"
)

func defaultTagger() *tagger.Tagger {
	// registers port 8080 as "http-alt"
	return tagger.New(map[uint16]string{8080: "http-alt"})
}

func TestPrivilegedPortIsCritical(t *testing.T) {
	c := severity.New(1023, nil)
	if got := c.Classify(80); got != severity.Critical {
		t.Fatalf("expected Critical for port 80, got %s", got)
	}
}

func TestBoundaryPortIsCritical(t *testing.T) {
	c := severity.New(1023, nil)
	if got := c.Classify(1023); got != severity.Critical {
		t.Fatalf("expected Critical for port 1023, got %s", got)
	}
}

func TestJustAbovePrivilegedWithNoTagIsWarning(t *testing.T) {
	c := severity.New(1023, defaultTagger())
	if got := c.Classify(9999); got != severity.Warning {
		t.Fatalf("expected Warning for port 9999, got %s", got)
	}
}

func TestTaggedHighPortIsInfo(t *testing.T) {
	c := severity.New(1023, defaultTagger())
	if got := c.Classify(8080); got != severity.Info {
		t.Fatalf("expected Info for tagged port 8080, got %s", got)
	}
}

func TestNilTaggerHighPortIsWarning(t *testing.T) {
	c := severity.New(1023, nil)
	if got := c.Classify(8080); got != severity.Warning {
		t.Fatalf("expected Warning when tagger is nil, got %s", got)
	}
}

func TestLevelStrings(t *testing.T) {
	cases := []struct {
		level severity.Level
		want  string
	}{
		{severity.Info, "INFO"},
		{severity.Warning, "WARNING"},
		{severity.Critical, "CRITICAL"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("Level(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}
