package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeState(proto, addr string, port uint16) scanner.PortState {
	return scanner.PortState{
		Protocol:    proto,
		LocalAddr:   addr,
		LocalPort:   port,
		ProcessName: "testproc",
	}
}

func TestNewAlerterDefaultsToStdout(t *testing.T) {
	a := alert.NewAlerter(nil)
	if a == nil {
		t.Fatal("expected non-nil Alerter")
	}
}

func TestNotifyWritesAlert(t *testing.T) {
	var buf bytes.Buffer
	a := alert.NewAlerter(&buf)
	state := makeState("tcp", "0.0.0.0", 8080)

	al := a.Notify(alert.LevelAlert, state)

	if al.Level != alert.LevelAlert {
		t.Errorf("expected level ALERT, got %s", al.Level)
	}
	if al.State.LocalPort != 8080 {
		t.Errorf("expected port 8080, got %d", al.State.LocalPort)
	}
	output := buf.String()
	if !strings.Contains(output, "ALERT") {
		t.Errorf("expected ALERT in output, got: %s", output)
	}
	if !strings.Contains(output, "unexpected listener") {
		t.Errorf("expected 'unexpected listener' in output, got: %s", output)
	}
}

func TestNotifyGoneWritesInfo(t *testing.T) {
	var buf bytes.Buffer
	a := alert.NewAlerter(&buf)
	state := makeState("tcp", "127.0.0.1", 9090)

	al := a.NotifyGone(state)

	if al.Level != alert.LevelInfo {
		t.Errorf("expected level INFO, got %s", al.Level)
	}
	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("expected INFO in output, got: %s", output)
	}
	if !strings.Contains(output, "listener closed") {
		t.Errorf("expected 'listener closed' in output, got: %s", output)
	}
}

func TestAlertTimestampSet(t *testing.T) {
	var buf bytes.Buffer
	a := alert.NewAlerter(&buf)
	state := makeState("udp", "0.0.0.0", 53)

	al := a.Notify(alert.LevelWarn, state)

	if al.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp on alert")
	}
}
