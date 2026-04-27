package portmap_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portmap"
	"github.com/user/portwatch/internal/scanner"
)

func makeState(port uint16, proto, addr string) scanner.PortState {
	return scanner.PortState{
		Port:     port,
		Proto:    proto,
		LocalIP:  addr,
		PID:      1234,
		Inode:    99,
	}
}

func TestSyncerAddsNewState(t *testing.T) {
	m := portmap.New()
	syncer := portmap.NewSyncer(m)

	st := makeState(8080, "tcp", "0.0.0.0")
	syncer.Add(st)

	key := portmap.KeyFromState(st)
	got, ok := m.Get(key)
	if !ok {
		t.Fatal("expected state to be present after Add")
	}
	if got.Port != st.Port {
		t.Errorf("port mismatch: got %d, want %d", got.Port, st.Port)
	}
}

func TestSyncerRemovesState(t *testing.T) {
	m := portmap.New()
	syncer := portmap.NewSyncer(m)

	st := makeState(9090, "udp", "127.0.0.1")
	syncer.Add(st)
	syncer.Remove(st)

	key := portmap.KeyFromState(st)
	_, ok := m.Get(key)
	if ok {
		t.Fatal("expected state to be absent after Remove")
	}
}

func TestSyncerSyncReplacesAll(t *testing.T) {
	m := portmap.New()
	syncer := portmap.NewSyncer(m)

	old := makeState(3000, "tcp", "0.0.0.0")
	syncer.Add(old)

	newStates := []scanner.PortState{
		makeState(4000, "tcp", "0.0.0.0"),
		makeState(5000, "udp", "0.0.0.0"),
	}
	syncer.Sync(newStates)

	// old key must be gone
	oldKey := portmap.KeyFromState(old)
	_, ok := m.Get(oldKey)
	if ok {
		t.Error("expected old state to be removed after Sync")
	}

	// new keys must be present
	for _, ns := range newStates {
		k := portmap.KeyFromState(ns)
		_, ok := m.Get(k)
		if !ok {
			t.Errorf("expected state for port %d to be present after Sync", ns.Port)
		}
	}
}

func TestSyncerSyncEmptySliceClearsMap(t *testing.T) {
	m := portmap.New()
	syncer := portmap.NewSyncer(m)

	syncer.Add(makeState(7777, "tcp", "0.0.0.0"))
	syncer.Sync(nil)

	if m.Len() != 0 {
		t.Errorf("expected map to be empty after Sync(nil), got len=%d", m.Len())
	}
}

func TestSyncerAddUpdatesTimestamp(t *testing.T) {
	m := portmap.New()
	syncer := portmap.NewSyncer(m)

	st := makeState(6443, "tcp", "0.0.0.0")
	before := time.Now()
	syncer.Add(st)
	after := time.Now()

	key := portmap.KeyFromState(st)
	entry, ok := m.GetEntry(key)
	if !ok {
		t.Fatal("entry not found")
	}
	if entry.LastSeen.Before(before) || entry.LastSeen.After(after) {
		t.Errorf("LastSeen %v not within [%v, %v]", entry.LastSeen, before, after)
	}
}

func TestSyncerDeltaReturnsAddedAndRemoved(t *testing.T) {
	m := portmap.New()
	syncer := portmap.NewSyncer(m)

	initial := []scanner.PortState{
		makeState(80, "tcp", "0.0.0.0"),
		makeState(443, "tcp", "0.0.0.0"),
	}
	syncer.Sync(initial)

	updated := []scanner.PortState{
		makeState(443, "tcp", "0.0.0.0"), // kept
		makeState(8080, "tcp", "0.0.0.0"), // new
	}
	added, removed := syncer.Delta(updated)

	if len(added) != 1 || added[0].Port != 8080 {
		t.Errorf("expected added=[8080], got %v", added)
	}
	if len(removed) != 1 || removed[0].Port != 80 {
		t.Errorf("expected removed=[80], got %v", removed)
	}
}
