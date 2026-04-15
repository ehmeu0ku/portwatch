// Package process provides utilities for resolving the owning process of a
// listening network socket on Linux.
//
// Given a socket inode number (obtained from /proc/net/tcp or /proc/net/tcp6),
// the Resolver walks /proc/<pid>/fd to find which process holds an open file
// descriptor pointing to that socket, then reads the process name from
// /proc/<pid>/comm.
//
// Usage:
//
//	r := process.NewResolver("/proc")
//	info, err := r.LookupInode(12345)
//	if err == nil {
//		fmt.Println(info) // pid=1234 name=nginx user=...
//	}
//
// Note: On non-Linux systems the resolver will return errors for every lookup
// because /proc is not available. Build-tag guards can be added if cross-
// platform support is required in the future.
package process
