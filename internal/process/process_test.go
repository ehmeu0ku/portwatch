package process_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/yourorg/portwatch/internal/process"
)

// buildFakeProc creates a minimal /proc-like directory tree under a temp dir.
// It creates one process entry with a single fd symlink pointing to the given
// socket inode, plus a comm file with the process name.
func buildFakeProc(t *testing.T, pid int, inode uint64, comm string) string {
	t.Helper()
	root := t.TempDir()

	pidDir := filepath.Join(root, strconv.Itoa(pid))
	fdDir := filepath.Join(pidDir, "fd")
	if err := os.MkdirAll(fdDir, 0o755); err != nil {
		t.Fatalf("mkdir fd: %v", err)
	}

	// comm file
	if err := os.WriteFile(filepath.Join(pidDir, "comm"), []byte(comm+"\n"), 0o644); err != nil {
		t.Fatalf("write comm: %v", err)
	}

	// fd/0 -> socket:[<inode>]
	target := fmt.Sprintf("socket:[%d]", inode)
	if err := os.Symlink(target, filepath.Join(fdDir, "0")); err != nil {
		t.Fatalf("symlink: %v", err)
	}

	return root
}

func TestLookupInodeFindsProcess(t *testing.T) {
	const pid = 4242
	const inode = 99887
	const comm = "nginx"

	root := buildFakeProc(t, pid, inode, comm)
	r := process.NewResolver(root)

	info, err := r.LookupInode(inode)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.PID != pid {
		t.Errorf("PID: got %d, want %d", info.PID, pid)
	}
	if info.Name != comm {
		t.Errorf("Name: got %q, want %q", info.Name, comm)
	}
}

func TestLookupInodeNotFound(t *testing.T) {
	root := buildFakeProc(t, 1, 111, "other")
	r := process.NewResolver(root)

	_, err := r.LookupInode(999) // inode that doesn't exist
	if err == nil {
		t.Fatal("expected error for missing inode, got nil")
	}
}

func TestInfoString(t *testing.T) {
	info := process.Info{PID: 7, Name: "sshd", User: "root"}
	got := info.String()
	want := "pid=7 name=sshd user=root"
	if got != want {
		t.Errorf("String(): got %q, want %q", got, want)
	}
}

func TestLookupInodeBadProcRoot(t *testing.T) {
	r := process.NewResolver("/nonexistent/proc/root")
	_, err := r.LookupInode(1)
	if err == nil {
		t.Fatal("expected error for bad proc root")
	}
}
