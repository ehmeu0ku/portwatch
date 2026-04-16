// Package webhook delivers alert events to an HTTP endpoint.
package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/correlator"
)

// Config holds webhook delivery settings.
type Config struct {
	URL     string
	Timeout time.Duration
	Secret  string // added as X-Portwatch-Secret header when non-empty
}

// Notifier sends correlator events to an HTTP endpoint as JSON.
type Notifier struct {
	cg    Config
	client *http.Client
}

// New returns a Notifier using.
func New(cfg Config) *Notifier {
	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}
	return &Notifier{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.Timeout},
	}
}

type payload struct {
	Kind      string `json:"kind"`
	Severity  string `json:"severity"`
	Port      uint16 `json:"port"`
	Proto     string `json:"proto"`
	Tag       string `json:"tag,omitempty"`
	Process   string `json:"process,omitempty"`
	Timestamp string `json:"timestamp"`
}

// Send delivers the event to the configured webhook URL.
func (n *Notifier) Send(ctx context.Context, ev correlator.Event) error {
	p := payload{
		Kind:      string(ev.Kind),
		Severity:  ev.Severity.String(),
		Port:      ev.State.Port,
		Proto:     ev.State.Proto,
		Tag:       ev.Tag,
		Timestamp: ev.Timestamp.Format(time.RFC3339),
	}
	if ev.State.Process != nil {
		p.Process = ev.State.Process.String()
	}
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("webhook: marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.cfg.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if n.cfg.Secret != "" {
		req.Header.Set("X-Portwatch-Secret", n.cfg.Secret)
	}
	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: send: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}

// Name implements notify.Backend.
func (n *Notifier) Name() string { return "webhook" }
