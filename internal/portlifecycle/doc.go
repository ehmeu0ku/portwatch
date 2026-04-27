// Package portlifecycle tracks the full lifecycle of observed ports.
//
// It records when a port was first seen, when it was last seen, how many
// times it has been observed, and when it was last marked gone. This
// information can be used by higher-level components to make decisions
// based on port stability, age, or churn rate.
//
// Usage:
//
//	tracker := portlifecycle.New()
//	entry := tracker.Observe("tcp:8080")
//	fmt.Println(entry.SeenCount) // 1
//	tracker.MarkGone("tcp:8080")
//	e, ok := tracker.Get("tcp:8080")
//	if ok && e.GoneAt != nil {
//		// port has been marked gone
//	}
package portlifecycle
