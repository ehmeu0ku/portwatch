package webhook_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/correlator"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/severity"
	"github.com/user/portwatch/internal/webhook"
)

func makeEvent(kind correlator.Kind) correlator.Event {
	return correlator.Event{
		Kind:      kind,
		Severity:  severity.Warning,
		Tag:       "http",
		Timestamp: time.Now(),
		State: scanner.PortState{
			Port:  8080,
			Proto: "tcp",
		},
	}
}

func TestSendPostsJSON(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("unexpected Content-Type: %s", ct)
		}
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	wh := webhook.New(webhook.Config{URL: srv.URL})
	if err := wh.Send(context.Background(), makeEvent(correlator.KindNew)); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if got["port"].(float64) != 8080 {
		t.Errorf("expected port 8080, got %v", got["port"])
	}
	if got["kind"] != "new" {
		t.Errorf("expected kind=new, got %v", got["kind"])
	}
}

func TestSendAttachesSecret(t *testing.T) {
	var secret string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret = r.Header.Get("X-Portwatch-Secret")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	wh := webhook.New(webhook.Config{URL: srv.URL, Secret: "tok3n"})
	if err := wh.Send(context.Background(), makeEvent(correlator.KindNew)); err != nil {
		t.Fatalf("Send: %v", err)
	}
	if secret != "tok3n" {
		t.Errorf("expected secret tok3n, got %q", secret)
	}
}

func TestSendReturnsErrorOnBadStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	wh := webhook.New(webhook.Config{URL: srv.URL})
	if err := wh.Send(context.Background(), makeEvent(correlator.KindNew)); err == nil {
		t.Error("expected error for 500 response")
	}
}

func TestNameIsWebhook(t *testing.T) {
	wh := webhook.New(webhook.Config{URL: "http://localhost"})
	if wh.Name() != "webhook" {
		t.Errorf("unexpected name: %s", wh.Name())
	}
}
