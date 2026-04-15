// Package config handles loading and validating portwatch configuration.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds the runtime configuration for portwatch.
type Config struct {
	// Interval is how often the scanner polls for port changes.
	Interval time.Duration `json:"interval"`

	// IgnorePorts is a list of ports to silently ignore.
	IgnorePorts []uint16 `json:"ignore_ports"`

	// AlertLogPath is the file path for alert output. Empty means stdout.
	AlertLogPath string `json:"alert_log_path"`

	// OnlyListenOn restricts monitoring to specific interfaces (e.g. "eth0").
	OnlyListenOn []string `json:"only_listen_on"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Interval:     5 * time.Second,
		IgnorePorts:  []uint16{},
		AlertLogPath: "",
		OnlyListenOn: []string{},
	}
}

// Load reads a JSON config file from path and merges it with defaults.
// If path is empty, the default config is returned.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()
	if path == "" {
		return cfg, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Validate checks that config values are within acceptable ranges.
func (c *Config) Validate() error {
	if c.Interval < time.Second {
		return fmt.Errorf("config: interval %v is too short (minimum 1s)", c.Interval)
	}
	return nil
}

// IsIgnored reports whether the given port should be suppressed from alerts.
func (c *Config) IsIgnored(port uint16) bool {
	for _, p := range c.IgnorePorts {
		if p == port {
			return true
		}
	}
	return false
}
