package debounce_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/debounce"
)

func TestTriggerFiresAfterWindow(t *testing.T) {
	var mu sync.Mutex
	fired := []string{}

	d := debounce.New(50*time.Millisecond, func(key string) {
		mu.Lock()
		fired = append(fired, key)
		mu.Unlock()
	})

	d.Trigger("tcp:8080")
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(fired) != 1 || fired[0] != "tcp:8080" {
		t.Fatalf("expected [tcp:8080], got %v", fired)
	}
}

func TestTriggerResetsTimer(t *testing.T) {
	var mu sync.Mutex
	count := 0

	d := debounce.New(80*time.Millisecond, func(_ string) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	// Trigger three times within the window; should only fire once.
	d.Trigger("key")
	time.Sleep(30 * time.Millisecond)
	d.Trigger("key")
	time.Sleep(30 * time.Millisecond)
	d.Trigger("key")
	time.Sleep(120 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if count != 1 {
		t.Fatalf("expected action to fire once, got %d", count)
	}
}

func TestCancelPreventsAction(t *testing.T) {
	fired := false

	d := debounce.New(60*time.Millisecond, func(_ string) {
		fired = true
	})

	d.Trigger("tcp:9090")
	d.Cancel("tcp:9090")
	time.Sleep(100 * time.Millisecond)

	if fired {
		t.Fatal("action should not have fired after Cancel")
	}
}

func TestPendingCount(t *testing.T) {
	d := debounce.New(200*time.Millisecond, func(_ string) {})

	if d.Pending() != 0 {
		t.Fatal("expected 0 pending timers initially")
	}

	d.Trigger("a")
	d.Trigger("b")

	if got := d.Pending(); got != 2 {
		t.Fatalf("expected 2 pending timers, got %d", got)
	}

	d.Cancel("a")
	if got := d.Pending(); got != 1 {
		t.Fatalf("expected 1 pending timer after cancel, got %d", got)
	}
}

func TestIndependentKeys(t *testing.T) {
	var mu sync.Mutex
	fired := map[string]int{}

	d := debounce.New(50*time.Millisecond, func(key string) {
		mu.Lock()
		fired[key]++
		mu.Unlock()
	})

	d.Trigger("x")
	d.Trigger("y")
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if fired["x"] != 1 || fired["y"] != 1 {
		t.Fatalf("expected each key to fire once, got %v", fired)
	}
}
