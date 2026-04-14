package scanner

import (
	"testing"
	"time"
)

func TestIsLoopback(t *testing.T) {
	cases := []struct {
		addr     string
		expected bool
	}{
		{"127.0.0.1", true},
		{"::1", true},
		{"0.0.0.0", false},
		{"192.168.1.1", false},
		{"invalid", false},
	}
	for _, tc := range cases {
		t.Run(tc.addr, func(t *testing.T) {
			got := IsLoopback(tc.addr)
			if got != tc.expected {
				t.Errorf("IsLoopback(%q) = %v, want %v", tc.addr, got, tc.expected)
			}
		})
	}
}

func TestDeduplicateStates(t *testing.T) {
	now := time.Now()
	input := []PortState{
		{Protocol: "tcp", Address: "0.0.0.0", Port: 8080, SeenAt: now},
		{Protocol: "tcp", Address: "0.0.0.0", Port: 8080, SeenAt: now},
		{Protocol: "tcp", Address: "0.0.0.0", Port: 443, SeenAt: now},
	}
	result := DeduplicateStates(input)
	if len(result) != 2 {
		t.Errorf("expected 2 unique states, got %d", len(result))
	}
}

func TestPortStateString(t *testing.T) {
	ps := PortState{
		Protocol: "tcp",
		Address:  "0.0.0.0",
		Port:     8080,
		PID:      1234,
		Process:  "nginx",
	}
	got := ps.String()
	expected := "tcp 0.0.0.0:8080 (pid=1234, process=nginx)"
	if got != expected {
		t.Errorf("String() = %q, want %q", got, expected)
	}
}

func TestProcScannerCreation(t *testing.T) {
	s := NewProcScanner()
	if s == nil {
		t.Fatal("NewProcScanner() returned nil")
	}
	// Scan should not panic even if /proc/net/tcp is unavailable
	_, err := s.Scan()
	if err != nil {
		t.Logf("Scan returned error (expected on non-Linux): %v", err)
	}
}
