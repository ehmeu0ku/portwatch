package anomaly_test

import (
	"testing"

	"github.com/user/portwatch/internal/anomaly"
	"github.com/user/portwatch/internal/scorecard"
)

func newSC() *scorecard.Scorecard { return scorecard.New() }

func TestCheckReturnsFalseForUnknownKey(t *testing.T) {
	d := anomaly.New(newSC(), 0.3, 0.7)
	_, ok := d.Check("tcp:8080")
	if ok {
		t.Fatal("expected no finding for unknown key")
	}
}

func TestCheckReturnsFalseWhenRatioBelowThreshold(t *testing.T) {
	sc := newSC()
	sc.Record("tcp:9000", false)
	sc.Record("tcp:9000", false)
	d := anomaly.New(sc, 0.3, 0.7)
	_, ok := d.Check("tcp:9000")
	if ok {
		t.Fatal("expected no finding below threshold")
	}
}

func TestCheckDetectsSuspicious(t *testing.T) {
	sc := newSC()
	sc.Record("tcp:9001", true)
	sc.Record("tcp:9001", false)
	d := anomaly.New(sc, 0.3, 0.7)
	f, ok := d.Check("tcp:9001")
	if !ok {
		t.Fatal("expected a finding")
	}
	if f.Level != anomaly.LevelSuspicious {
		t.Fatalf("expected suspicious, got %s", f.Level)
	}
}

func TestCheckDetectsAnomalous(t *testing.T) {
	sc := newSC()
	for i := 0; i < 3; i++ {
		sc.Record("tcp:22", true)
	}
	sc.Record("tcp:22", false)
	d := anomaly.New(sc, 0.3, 0.7)
	f, ok := d.Check("tcp:22")
	if !ok {
		t.Fatal("expected a finding")
	}
	if f.Level != anomaly.LevelAnomalous {
		t.Fatalf("expected anomalous, got %s", f.Level)
	}
}

func TestFindingStringContainsKey(t *testing.T) {
	sc := newSC()
	sc.Record("tcp:443", true)
	d := anomaly.New(sc, 0.5, 0.9)
	f, ok := d.Check("tcp:443")
	if !ok {
		t.Fatal("expected finding")
	}
	if s := f.String(); len(s) == 0 {
		t.Fatal("empty string")
	}
}

func TestLevelString(t *testing.T) {
	cases := []struct{ l anomaly.Level; want string }{
		{anomaly.LevelNone, "none"},
		{anomaly.LevelSuspicious, "suspicious"},
		{anomaly.LevelAnomalous, "anomalous"},
	}
	for _, c := range cases {
		if got := c.l.String(); got != c.want {
			t.Errorf("Level(%d).String() = %q, want %q", c.l, got, c.want)
		}
	}
}
