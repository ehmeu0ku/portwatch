// Package scanner provides interfaces and implementations for discovering
// active network port listeners on the local machine.
//
// The primary interface is Scanner, which any backend must implement:
//
//	type Scanner interface {
//	    Scan() ([]PortState, error)
//	}
//
// Implementations:
//   - ProcScanner: reads /proc/net/tcp and /proc/net/tcp6 (Linux only)
//   - MockScanner: in-memory scanner for use in tests
//
// PortState holds metadata about a single listening port, including protocol,
// address, port number, PID, and the process name when available.
//
// Usage:
//
//	s := scanner.NewProcScanner()
//	states, err := s.Scan()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, state := range states {
//	    fmt.Println(state)
//	}
package scanner
