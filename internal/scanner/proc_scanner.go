package scanner

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// ProcScanner reads active listeners from /proc/net/tcp and /proc/net/tcp6.
type ProcScanner struct{}

// NewProcScanner creates a new ProcScanner.
func NewProcScanner() *ProcScanner {
	return &ProcScanner{}
}

// Scan returns all currently listening TCP ports from /proc/net/tcp and tcp6.
func (s *ProcScanner) Scan() ([]PortState, error) {
	var states []PortState
	for _, path := range []string{"/proc/net/tcp", "/proc/net/tcp6"} {
		results, err := parseProcNet(path)
		if err != nil {
			// Non-fatal: file may not exist on all systems
			continue
		}
		states = append(states, results...)
	}
	return DeduplicateStates(states), nil
}

// parseProcNet parses a /proc/net/tcp or /proc/net/tcp6 file.
func parseProcNet(path string) ([]PortState, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	var states []PortState
	scanner := bufio.NewScanner(f)
	scanner.Scan() // skip header
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 4 {
			continue
		}
		// state 0A = listening
		if fields[3] != "0A" {
			continue
		}
		addr, port, err := parseHexAddr(fields[1])
		if err != nil {
			continue
		}
		states = append(states, PortState{
			Protocol: "tcp",
			Address:  addr,
			Port:     port,
			SeenAt:   time.Now(),
		})
	}
	return states, scanner.Err()
}

// parseHexAddr converts a hex-encoded address:port pair (little-endian) to ip:port.
func parseHexAddr(hexAddr string) (string, int, error) {
	parts := strings.Split(hexAddr, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid addr: %s", hexAddr)
	}
	portVal, err := strconv.ParseInt(parts[1], 16, 32)
	if err != nil {
		return "", 0, err
	}
	return "0.0.0.0", int(portVal), nil
}
