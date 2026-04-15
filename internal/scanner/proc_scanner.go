package scanner

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
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
// The address field in /proc/net/tcp is stored as a little-endian 32-bit hex value
// for IPv4, or four little-endian 32-bit hex values for IPv6.
func parseHexAddr(hexAddr string) (string, int, error) {
	parts := strings.Split(hexAddr, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid addr: %s", hexAddr)
	}
	portVal, err := strconv.ParseInt(parts[1], 16, 32)
	if err != nil {
		return "", 0, err
	}

	ip, err := parseHexIP(parts[0])
	if err != nil {
		// Fall back to generic address if parsing fails
		return "0.0.0.0", int(portVal), nil
	}
	return ip, int(portVal), nil
}

// parseHexIP decodes a little-endian hex-encoded IP address from /proc/net/tcp.
// IPv4 addresses are 8 hex chars; IPv6 addresses are 32 hex chars.
func parseHexIP(hexIP string) (string, error) {
	b, err := hex.DecodeString(hexIP)
	if err != nil {
		return "", fmt.Errorf("decode hex IP %q: %w", hexIP, err)
	}
	switch len(b) {
	case 4:
		// IPv4: reverse byte order (little-endian)
		v := binary.LittleEndian.Uint32(b)
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, v)
		return ip.String(), nil
	case 16:
		// IPv6: four little-endian 32-bit words
		ip := make(net.IP, 16)
		for i := 0; i < 4; i++ {
			v := binary.LittleEndian.Uint32(b[i*4 : i*4+4])
			binary.BigEndian.PutUint32(ip[i*4:], v)
		}
		return ip.String(), nil
	default:
		return "", fmt.Errorf("unexpected IP length %d for %q", len(b), hexIP)
	}
}
