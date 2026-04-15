// Package reporter provides formatted summary reporting of port monitoring activity.
package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Reporter writes periodic summaries of port state to an output writer.
type Reporter struct {
	out       io.Writer
	startTime time.Time
}

// New creates a new Reporter writing to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{
		out:       w,
		startTime: time.Now(),
	}
}

// Summary writes a formatted summary of the current port states to the output.
func (r *Reporter) Summary(states []scanner.PortState) {
	fmt.Fprintf(r.out, "--- portwatch summary at %s ---\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(r.out, "uptime: %s\n", time.Since(r.startTime).Round(time.Second))

	if len(states) == 0 {
		fmt.Fprintln(r.out, "no active listeners detected")
		return
	}

	sorted := make([]scanner.PortState, len(states))
	copy(sorted, states)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Proto != sorted[j].Proto {
			return sorted[i].Proto < sorted[j].Proto
		}
		return sorted[i].Port < sorted[j].Port
	})

	fmt.Fprintf(r.out, "active listeners (%d):\n", len(sorted))
	for _, s := range sorted {
		pid := "-"
		if s.PID > 0 {
			pid = fmt.Sprintf("%d", s.PID)
		}
		fmt.Fprintf(r.out, "  %-5s %s:%-5d pid=%-6s %s\n",
			s.Proto, s.IP, s.Port, pid, s.Process)
	}
}

// ReportNew writes a concise line indicating a newly detected port.
func (r *Reporter) ReportNew(s scanner.PortState) {
	fmt.Fprintf(r.out, "[NEW]  %s %s:%d (pid=%d %s) at %s\n",
		s.Proto, s.IP, s.Port, s.PID, s.Process, time.Now().Format(time.RFC3339))
}

// ReportGone writes a concise line indicating a port that stopped listening.
func (r *Reporter) ReportGone(s scanner.PortState) {
	fmt.Fprintf(r.out, "[GONE] %s %s:%d (pid=%d %s) at %s\n",
		s.Proto, s.IP, s.Port, s.PID, s.Process, time.Now().Format(time.RFC3339))
}
