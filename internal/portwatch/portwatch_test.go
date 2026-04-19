package portwatch_test

import (
	"testing"

	"github.com/user/portwatch/internal/portwatch"
	"github.com/user/portwatch/internal/process"
	"github.com/user/portwatch/internal/scanner"
)

func makeState(port uint16, proto string, pid int, comm string) scanner.PortState {
	return scanner.PortState{
		Port:    port,
		Proto:   proto,
		Process: &process.Info{PID: pid, Comm: comm},
	}
}

func TestClaimOnFirstSeen(t *testing.T) {
	tr := portwatch.New()
	changes := tr.Update([]scanner.PortState{makeState(8080, "tcp", 100, "nginx")})
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != portwatch.OwnershipClaimed {
		t.Errorf("expected claimed, got %s", changes[0].Kind)
	}
}

func TestNoChangeWhenSamePID(t *testing.T) {
	tr := portwatch.New()
	state := []scanner.PortState{makeState(8080, "tcp", 100, "nginx")}
	tr.Update(state)
	changes := tr.Update(state)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestOwnershipChangedOnPIDSwitch(t *testing.T) {
	tr := portwatch.New()
	tr.Update([]scanner.PortState{makeState(9090, "tcp", 1, "old")})
	changes := tr.Update([]scanner.PortState{makeState(9090, "tcp", 2, "new")})
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != portwatch.OwnershipChanged {
		t.Errorf("expected changed, got %s", changes[0].Kind)
	}
	if changes[0].Prev.PID != 1 || changes[0].Current.PID != 2 {
		t.Errorf("unexpected PIDs: prev=%d current=%d", changes[0].Prev.PID, changes[0].Current.PID)
	}
}

func TestOwnershipReleasedWhenPortGone(t *testing.T) {
	tr := portwatch.New()
	tr.Update([]scanner.PortState{makeState(7070, "tcp", 5, "svc")})
	changes := tr.Update([]scanner.PortState{})
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != portwatch.OwnershipReleased {
		t.Errorf("expected released, got %s", changes[0].Kind)
	}
}

func TestChangeString(t *testing.T) {
	e := &portwatch.Entry{Port: 80, Proto: "tcp", PID: 10, Comm: "apache"}
	c := portwatch.Change{Kind: portwatch.OwnershipClaimed, Current: e}
	if c.String() == "" {
		t.Error("expected non-empty string")
	}
}

func TestSkipsStatesWithNoProcess(t *testing.T) {
	tr := portwatch.New()
	changes := tr.Update([]scanner.PortState{{Port: 1234, Proto: "tcp", Process: nil}})
	if len(changes) != 0 {
		t.Errorf("expected 0 changes for nil process, got %d", len(changes))
	}
}
