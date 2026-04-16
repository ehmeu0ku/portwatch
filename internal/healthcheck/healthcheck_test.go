package healthcheck

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleHealthReturnsOK(t *testing.T) {
	s := New(":0")
	rec := httptest.NewRecorder()
	s.handleHealth(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var st Status
	if err := json.NewDecoder(rec.Body).Decode(&st); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !st.OK {
		t.Error("expected ok=true")
	}
}

func TestRecordScanIncrements(t *testing.T) {
	s := New(":0")
	s.RecordScan()
	s.RecordScan()
	if got := s.scans.Load(); got != 2 {
		t.Fatalf("expected 2 scans, got %d", got)
	}
}

func TestStartedAtIsSet(t *testing.T) {
	before := time.Now()
	s := New(":0")
	if s.startedAt.Before(before) {
		t.Error("startedAt should be >= before")
	}
}

func TestStartAndShutdown(t *testing.T) {
	s := New("127.0.0.1:0")
	// Use a real port via httptest to avoid binding conflicts in CI.
	srv := httptest.NewServer(http.HandlerFunc(s.handleHealth))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = ctx

	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestScanCountAppearsInResponse(t *testing.T) {
	s := New(":0")
	s.RecordScan()
	s.RecordScan()
	s.RecordScan()
	rec := httptest.NewRecorder()
	s.handleHealth(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	var st Status
	_ = json.NewDecoder(rec.Body).Decode(&st)
	if st.Scans != 3 {
		t.Fatalf("expected 3, got %d", st.Scans)
	}
}
