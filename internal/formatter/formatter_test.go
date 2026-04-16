package formatter_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/correlator"
	"github.com/user/portwatch/internal/formatter"
	"github.com/user/portwatch/internal/process"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/severity"
)

func makeEvent(port uint16, kind string, sev severity.Level) correlator.Event {
	return correlator.Event{
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Kind:      kind,
		Severity:  sev,
		Tag:       "http",
		State: scanner.PortState{
			Proto: "tcp",
			IP:    "0.0.0.0",
			Port:  port,
			PID:   1234,
			Process: &process.Info{
				Name: "nginx",
				PID:  1234,
			},
		},
	}
}

func TestTextFormatContainsKind(t *testing.T) {
	f := formatter.New(formatter.StyleText)
	out := f.Format(makeEvent(80, "NEW", severity.Critical))
	if !strings.Contains(out, "NEW") {
		t.Errorf("expected NEW in output, got: %s", out)
	}
}

func TestTextFormatContainsSeverity(t *testing.T) {
	f := formatter.New(formatter.StyleText)
	out := f.Format(makeEvent(80, "NEW", severity.Critical))
	if !strings.Contains(out, "CRITICAL") {
		t.Errorf("expected CRITICAL in output, got: %s", out)
	}
}

func TestTextFormatContainsPort(t *testing.T) {
	f := formatter.New(formatter.StyleText)
	out := f.Format(makeEvent(8080, "NEW", severity.Warning))
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", out)
	}
}

func TestTextFormatContainsProcess(t *testing.T) {
	f := formatter.New(formatter.StyleText)
	out := f.Format(makeEvent(80, "NEW", severity.Info))
	if !strings.Contains(out, "nginx") {
		t.Errorf("expected process name in output, got: %s", out)
	}
}

func TestJSONFormatIsValidJSON(t *testing.T) {
	f := formatter.New(formatter.StyleJSON)
	out := f.Format(makeEvent(443, "NEW", severity.Critical))
	if !strings.HasPrefix(out, "{") || !strings.HasSuffix(out, "}") {
		t.Errorf("expected JSON object, got: %s", out)
	}
}

func TestJSONFormatContainsFields(t *testing.T) {
	f := formatter.New(formatter.StyleJSON)
	out := f.Format(makeEvent(443, "GONE", severity.Warning))
	for _, field := range []string{"ts", "kind", "severity", "proto", "ip", "port", "pid", "process", "tag"} {
		if !strings.Contains(out, `"`+field+`"`) {
			t.Errorf("expected field %q in JSON output: %s", field, out)
		}
	}
}

func TestNilProcessRenderedAsDash(t *testing.T) {
	f := formatter.New(formatter.StyleText)
	e := makeEvent(22, "NEW", severity.Info)
	e.State.Process = nil
	out := f.Format(e)
	if !strings.Contains(out, "proc=-") {
		t.Errorf("expected proc=- for nil process, got: %s", out)
	}
}
