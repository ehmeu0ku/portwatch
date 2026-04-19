// Package portlock promotes ports to "locked" status after they have
// been observed continuously for a configurable minimum age.
//
// A locked port is treated as an expected permanent listener; alerting
// pipelines can use this to suppress repetitive notifications for
// services like sshd or nginx that are always present.
//
// Typical usage:
//
//	store := portlock.New(5 * time.Minute)
//	if store.Observe("tcp:22", time.Now()) {
//		log.Println("tcp:22 is now locked")
//	}
package portlock
