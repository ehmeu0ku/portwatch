package reporter_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/reporter"
	"github.com/user/portwatch/internal/scanner"
)

func makeState(proto, ip string, port int, pid int, process string) scanner.PortState {
	return scanner.PortState{
		Proto:   proto,
		IP:      ip,
		Port:    port,
		PID:     pid,
		Process: process,
	}
}

func TestSummaryEmpty(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf)
	r.Summary(nil)
	if !strings.Contains(buf.String(), "no active listeners") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestSummaryListsStates(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf)
	states := []scanner.PortState{
		makeState("tcp", "0.0.0.0", 8080, 1234, "myapp"),
		makeState("udp", "127.0.0.1", 53, 99, "dnsmasq"),
	}
	r.Summary(states)
	out := buf.String()
	if !strings.Contains(out, "active listeners (2)") {
		t.Errorf("expected listener count, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output")
	}
	if !strings.Contains(out, "dnsmasq") {
		t.Errorf("expected process name in output")
	}
}

func TestSummaryIsSorted(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf)
	states := []scanner.PortState{
		makeState("tcp", "0.0.0.0", 9000, 1, "z"),
		makeState("tcp", "0.0.0.0", 80, 2, "a"),
	}
	r.Summary(states)
	out := buf.String()
	idx80 := strings.Index(out, ":80")
	idx9000 := strings.Index(out, ":9000")
	if idx80 > idx9000 {
		t.Errorf("expected port 80 before 9000 in sorted output")
	}
}

func TestReportNewContainsNEW(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf)
	r.ReportNew(makeState("tcp", "0.0.0.0", 443, 555, "nginx"))
	if !strings.Contains(buf.String(), "[NEW]") {
		t.Errorf("expected [NEW] tag in output")
	}
}

func TestReportGoneContainsGONE(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf)
	r.ReportGone(makeState("tcp", "0.0.0.0", 8080, 100, "server"))
	if !strings.Contains(buf.String(), "[GONE]") {
		t.Errorf("expected [GONE] tag in output")
	}
}

func TestNewDefaultsToStdout(t *testing.T) {
	// Just ensure New(nil) does not panic
	r := reporter.New(nil)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}

func TestSummaryContainsTimestamp(t *testing.T) {
	var buf strings.Builder
	r := reporter.New(&buf)
	r.Summary(nil)
	year := time.Now().Format("2006")
	if !strings.Contains(buf.String(), year) {
		t.Errorf("expected current year in summary timestamp")
	}
}
