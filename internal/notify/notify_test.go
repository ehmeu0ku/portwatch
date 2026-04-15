package notify_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/notify"
)

var fixedTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func makeMsg(level notify.Level, title, body string) notify.Message {
	return notify.Message{Level: level, Title: title, Body: body, Timestamp: fixedTime}
}

func TestLogNotifierWritesFormattedLine(t *testing.T) {
	var buf bytes.Buffer
	n := notify.NewLogNotifier(&buf)

	if err := n.Send(makeMsg(notify.LevelAlert, "new port", "TCP :9090")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	for _, want := range []string{"ALERT", "new port", "TCP :9090", "2024-06-01T12:00:00Z"} {
		if !strings.Contains(got, want) {
			t.Errorf("output %q missing %q", got, want)
		}
	}
}

func TestLogNotifierName(t *testing.T) {
	n := notify.NewLogNotifier(nil)
	if n.Name() != "log" {
		t.Errorf("expected 'log', got %q", n.Name())
	}
}

func TestMultiNotifierDeliversToBothBackends(t *testing.T) {
	var b1, b2 bytes.Buffer
	multi := notify.NewMultiNotifier(notify.NewLogNotifier(&b1), notify.NewLogNotifier(&b2))

	if err := multi.Send(makeMsg(notify.LevelWarn, "gone", "UDP :53")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, buf := range []*bytes.Buffer{&b1, &b2} {
		if !strings.Contains(buf.String(), "UDP :53") {
			t.Errorf("backend %d did not receive message", i+1)
		}
	}
}

func TestMultiNotifierCollectsErrors(t *testing.T) {
	failing := &errorNotifier{}
	multi := notify.NewMultiNotifier(failing, failing)

	err := multi.Send(makeMsg(notify.LevelInfo, "x", "y"))
	if err == nil {
		t.Fatal("expected error from failing backends")
	}
	if !strings.Contains(err.Error(), "2 backend(s) failed") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestMultiNotifierAdd(t *testing.T) {
	var buf bytes.Buffer
	multi := notify.NewMultiNotifier()
	multi.Add(notify.NewLogNotifier(&buf))

	_ = multi.Send(makeMsg(notify.LevelInfo, "t", "b"))
	if buf.Len() == 0 {
		t.Error("expected output after Add")
	}
}

// errorNotifier is a test helper that always returns an error.
type errorNotifier struct{}

func (e *errorNotifier) Name() string          { return "error" }
func (e *errorNotifier) Send(_ notify.Message) error { return errors.New("send failed") }
