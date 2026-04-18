package tunables_test

import (
	"testing"
	"time"

	"github.com/example/portwatch/internal/tunables"
)

func TestDefaultsAreReasonable(t *testing.T) {
	tn := tunables.Defaults()
	if tn.ScanInterval() < time.Second {
		t.Fatal("scan interval too short")
	}
	if tn.AlertCooldown() <= 0 {
		t.Fatal("alert cooldown must be positive")
	}
	if tn.MaxHistorySize() < 1 {
		t.Fatal("max history size must be at least 1")
	}
}

func TestSetScanIntervalAcceptsValid(t *testing.T) {
	tn := tunables.Defaults()
	if !tn.SetScanInterval(10 * time.Second) {
		t.Fatal("expected true")
	}
	if tn.ScanInterval() != 10*time.Second {
		t.Fatalf("got %v", tn.ScanInterval())
	}
}

func TestSetScanIntervalRejectsTooShort(t *testing.T) {
	tn := tunables.Defaults()
	if tn.SetScanInterval(500 * time.Millisecond) {
		t.Fatal("expected false for sub-second interval")
	}
}

func TestSetAlertCooldownAcceptsZero(t *testing.T) {
	tn := tunables.Defaults()
	if !tn.SetAlertCooldown(0) {
		t.Fatal("zero cooldown should be allowed")
	}
	if tn.AlertCooldown() != 0 {
		t.Fatalf("got %v", tn.AlertCooldown())
	}
}

func TestSetAlertCooldownRejectsNegative(t *testing.T) {
	tn := tunables.Defaults()
	if tn.SetAlertCooldown(-1 * time.Second) {
		t.Fatal("expected false for negative cooldown")
	}
}

func TestSetMaxHistorySizeAcceptsPositive(t *testing.T) {
	tn := tunables.Defaults()
	if !tn.SetMaxHistorySize(100) {
		t.Fatal("expected true")
	}
	if tn.MaxHistorySize() != 100 {
		t.Fatalf("got %d", tn.MaxHistorySize())
	}
}

func TestSetMaxHistorySizeRejectsZero(t *testing.T) {
	tn := tunables.Defaults()
	if tn.SetMaxHistorySize(0) {
		t.Fatal("expected false for zero size")
	}
}
