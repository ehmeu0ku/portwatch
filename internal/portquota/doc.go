// Package portquota provides per-protocol listener count quotas for
// portwatch.
//
// A Quota instance tracks how many ports are currently active for each
// network protocol (e.g. "tcp", "udp") and compares that count against
// a configurable ceiling. When the ceiling is breached, Observe returns
// true so the caller can raise an alert or suppress further scanning.
//
// Usage:
//
//	q := portquota.New(50)          // default ceiling of 50
//	q.SetCeiling("tcp", 100)        // override for TCP
//
//	if q.Observe("tcp") {
//	    // ceiling just exceeded — alert
//	}
//	defer q.Release("tcp")
package portquota
