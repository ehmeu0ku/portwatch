package monitor

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func TestMonitorDetectsNewPort(t *testing.T) {
	var buf bytes.Buffer
	al := alert.NewAlerter(&buf)

	// Scanner returns a new port on the second call.
	callCount := 0
	sc := scanner.NewMockScanner(func() ([]scanner.PortState, error) {
		callCount++
		if callCount == 1 {
			return nil, nil
		}
		return []scanner.PortState{
			{Proto: "tcp", Addr: "0.0.0.0", Port: 4444},
		}, nil
	})

	cfg := Config{Interval: 20 * time.Millisecond}
	m := New(cfg, sc, al)

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	_ = m.Run(ctx) // returns context.DeadlineExceeded — expected

	output := buf.String()
	if !strings.Contains(output, "4444") {
		t.Fatalf("expected alert for port 4444, got: %q", output)
	}
}

func TestMonitorDetectsGonePort(t *testing.T) {
	var buf bytes.Buffer
	al := alert.NewAlerter(&buf)

	callCount := 0
	sc := scanner.NewMockScanner(func() ([]scanner.PortState, error) {
		callCount++
		if callCount == 1 {
			return []scanner.PortState{
				{Proto: "tcp", Addr: "0.0.0.0", Port: 9999},
			}, nil
		}
		return nil, nil
	})

	cfg := Config{Interval: 20 * time.Millisecond}
	m := New(cfg, sc, al)

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	_ = m.Run(ctx)

	output := buf.String()
	if !strings.Contains(output, "9999") {
		t.Fatalf("expected gone-alert for port 9999, got: %q", output)
	}
}
