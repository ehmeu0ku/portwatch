package porttrend

import (
	"testing"
)

func TestStableWhenOneSample(t *testing.T) {
	tr := New(6)
	tr.Record("tcp:8080", true)
	if got := tr.Trend("tcp:8080"); got != Stable {
		t.Fatalf("expected Stable, got %s", got)
	}
}

func TestRisingWhenMoreRecentlySeen(t *testing.T) {
	tr := New(4)
	tr.Record("tcp:9090", false)
	tr.Record("tcp:9090", false)
	tr.Record("tcp:9090", true)
	tr.Record("tcp:9090", true)
	if got := tr.Trend("tcp:9090"); got != Rising {
		t.Fatalf("expected Rising, got %s", got)
	}
}

func TestFallingWhenLessRecentlySeen(t *testing.T) {
	tr := New(4)
	tr.Record("tcp:443", true)
	tr.Record("tcp:443", true)
	tr.Record("tcp:443", false)
	tr.Record("tcp:443", false)
	if got := tr.Trend("tcp:443"); got != Falling {
		t.Fatalf("expected Falling, got %s", got)
	}
}

func TestStableWhenEqualHalves(t *testing.T) {
	tr := New(4)
	tr.Record("udp:53", true)
	tr.Record("udp:53", false)
	tr.Record("udp:53", true)
	tr.Record("udp:53", false)
	if got := tr.Trend("udp:53"); got != Stable {
		t.Fatalf("expected Stable, got %s", got)
	}
}

func TestWindowEvictsOldSamples(t *testing.T) {
	tr := New(4)
	// old samples: all false
	for i := 0; i < 10; i++ {
		tr.Record("tcp:22", false)
	}
	// recent samples: all true
	tr.Record("tcp:22", true)
	tr.Record("tcp:22", true)
	tr.Record("tcp:22", true)
	tr.Record("tcp:22", true)
	if got := tr.Trend("tcp:22"); got != Stable {
		// all four kept samples are true so both halves equal → Stable
		t.Fatalf("expected Stable after full window of true, got %s", got)
	}
}

func TestMissingKeyIsStable(t *testing.T) {
	tr := New(6)
	if got := tr.Trend("tcp:12345"); got != Stable {
		t.Fatalf("expected Stable for unknown key, got %s", got)
	}
}

func TestDirectionString(t *testing.T) {
	cases := []struct {
		d    Direction
		want string
	}{
		{Rising, "rising"},
		{Falling, "falling"},
		{Stable, "stable"},
	}
	for _, c := range cases {
		if got := c.d.String(); got != c.want {
			t.Errorf("Direction(%d).String() = %q, want %q", c.d, got, c.want)
		}
	}
}
