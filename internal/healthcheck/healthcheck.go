// Package healthcheck exposes a simple HTTP endpoint that reports daemon liveness.
package healthcheck

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

// Status holds the payload returned by the health endpoint.
type Status struct {
	OK        bool      `json:"ok"`
	Uptime    string    `json:"uptime"`
	Scans     uint64    `json:"scans"`
	StartedAt time.Time `json:"started_at"`
}

// Server is a lightweight HTTP health-check server.
type Server struct {
	addr      string
	scans     atomic.Uint64
	startedAt time.Time
	server    *http.Server
}

// New creates a Server that will listen on addr (e.g. ":9090").
func New(addr string) *Server {
	s := &Server{addr: addr, startedAt: time.Now()}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	s.server = &http.Server{Addr: addr, Handler: mux}
	return s
}

// RecordScan increments the scan counter; call after every scanner tick.
func (s *Server) RecordScan() { s.scans.Add(1) }

// Start begins serving in a background goroutine. It returns when the server
// is ready or the context is cancelled.
func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()
	select {
	case err := <-errCh:
		return fmt.Errorf("healthcheck: %w", err)
	case <-time.After(20 * time.Millisecond):
		go func() {
			<-ctx.Done()
			_ = s.server.Shutdown(context.Background())
		}()
		return nil
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(Status{
		OK:        true,
		Uptime:    time.Since(s.startedAt).Round(time.Second).String(),
		Scans:     s.scans.Load(),
		StartedAt: s.startedAt,
	})
}
