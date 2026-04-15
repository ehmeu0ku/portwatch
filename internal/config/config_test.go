package config_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.Interval != 5*time.Second {
		t.Errorf("expected default interval 5s, got %v", cfg.Interval)
	}
	if cfg.AlertLogPath != "" {
		t.Errorf("expected empty alert log path, got %q", cfg.AlertLogPath)
	}
}

func TestLoadEmptyPath(t *testing.T) {
	cfg, err := config.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
}

func TestLoadFromFile(t *testing.T) {
	data := map[string]interface{}{
		"interval":       "10s",
		"ignore_ports":   []int{22, 80},
		"alert_log_path": "/tmp/alerts.log",
	}
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(data); err != nil {
		t.Fatal(err)
	}
	f.Close()

	cfg, err := config.Load(f.Name())
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.Interval != 10*time.Second {
		t.Errorf("expected 10s interval, got %v", cfg.Interval)
	}
	if cfg.AlertLogPath != "/tmp/alerts.log" {
		t.Errorf("unexpected alert log path: %q", cfg.AlertLogPath)
	}
	if len(cfg.IgnorePorts) != 2 {
		t.Errorf("expected 2 ignored ports, got %d", len(cfg.IgnorePorts))
	}
}

func TestValidateRejectsTooShortInterval(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Interval = 500 * time.Millisecond
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for sub-second interval")
	}
}

func TestIsIgnored(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.IgnorePorts = []uint16{22, 443}

	if !cfg.IsIgnored(22) {
		t.Error("expected port 22 to be ignored")
	}
	if cfg.IsIgnored(8080) {
		t.Error("expected port 8080 not to be ignored")
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/path/portwatch.json")
	if err == nil {
		t.Error("expected error for missing config file")
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("{not valid json"); err != nil {
		t.Fatal(err)
	}
	f.Close()

	_, err = config.Load(f.Name())
	if err == nil {
		t.Error("expected error when loading file with invalid JSON")
	}
}
