package scanner

import "time"

// MockScanner is a test double that returns a fixed set of PortStates.
type MockScanner struct {
	States []PortState
	Err    error
	Calls  int
}

// NewMockScanner creates a MockScanner pre-loaded with the given states.
func NewMockScanner(states []PortState) *MockScanner {
	return &MockScanner{States: states}
}

// Scan returns the pre-configured states and increments the call counter.
func (m *MockScanner) Scan() ([]PortState, error) {
	m.Calls++
	if m.Err != nil {
		return nil, m.Err
	}
	return m.States, nil
}

// DefaultTestStates returns a small set of realistic PortStates for testing.
func DefaultTestStates() []PortState {
	now := time.Now()
	return []PortState{
		{Protocol: "tcp", Address: "0.0.0.0", Port: 22, PID: 100, Process: "sshd", SeenAt: now},
		{Protocol: "tcp", Address: "127.0.0.1", Port: 5432, PID: 200, Process: "postgres", SeenAt: now},
		{Protocol: "tcp", Address: "0.0.0.0", Port: 80, PID: 300, Process: "nginx", SeenAt: now},
	}
}
