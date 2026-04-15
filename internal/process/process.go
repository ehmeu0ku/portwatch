// Package process resolves the owning process (PID, name, user) for a
// listening port by inspecting /proc on Linux systems.
package process

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Info holds process metadata associated with a listening socket.
type Info struct {
	PID  int
	Name string
	User string
}

func (i Info) String() string {
	return fmt.Sprintf("pid=%d name=%s user=%s", i.PID, i.Name, i.User)
}

// Resolver looks up process information for a given inode.
type Resolver struct {
	procRoot string
}

// NewResolver returns a Resolver that reads from the given procRoot
// (pass "/proc" for production use).
func NewResolver(procRoot string) *Resolver {
	return &Resolver{procRoot: procRoot}
}

// LookupInode returns the Info for the process that owns the socket inode.
// Returns an error if no matching process is found.
func (r *Resolver) LookupInode(inode uint64) (Info, error) {
	target := fmt.Sprintf("socket:[%d]", inode)

	entries, err := os.ReadDir(r.procRoot)
	if err != nil {
		return Info{}, fmt.Errorf("read proc: %w", err)
	}

	for _, e := range entries {
		pid, err := strconv.Atoi(e.Name())
		if err != nil {
			continue // skip non-numeric entries
		}
		fdDir := filepath.Join(r.procRoot, e.Name(), "fd")
		links, err := os.ReadDir(fdDir)
		if err != nil {
			continue
		}
		for _, fd := range links {
			link, err := os.Readlink(filepath.Join(fdDir, fd.Name()))
			if err != nil {
				continue
			}
			if link == target {
				name := r.readComm(pid)
				user := r.readUser(pid)
				return Info{PID: pid, Name: name, User: user}, nil
			}
		}
	}
	return Info{}, fmt.Errorf("inode %d: no owning process found", inode)
}

func (r *Resolver) readComm(pid int) string {
	data, err := os.ReadFile(filepath.Join(r.procRoot, strconv.Itoa(pid), "comm"))
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(data))
}

func (r *Resolver) readUser(pid int) string {
	info, err := os.Stat(filepath.Join(r.procRoot, strconv.Itoa(pid)))
	if err != nil {
		return "unknown"
	}
	// Surface the numeric UID via the FileInfo string representation.
	// A real implementation would look up /etc/passwd; this keeps it portable.
	return fmt.Sprintf("%v", info.Sys())
}
