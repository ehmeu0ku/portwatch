// Package watchdog implements a heartbeat-based liveness monitor for the
// portwatch scan loop.
//
// # Overview
//
// The Watchdog expects periodic Beat() calls from the component it supervises.
// If no beat is received within the configured timeout window, the registered
// onTimeout callback is invoked in a separate goroutine so that the caller can
// attempt recovery (e.g. restart the scanner, emit an alert, or exit).
//
// # Usage
//
//	wd := watchdog.New(30*time.Second, func() {
//		log.Println("scanner stalled — restarting")
//		// recovery logic here
//	})
//
//	go wd.Start(ctx)
//
//	// inside the scan loop:
//	wd.Beat()
//
// The watchdog stops automatically when the provided context is cancelled.
package watchdog
