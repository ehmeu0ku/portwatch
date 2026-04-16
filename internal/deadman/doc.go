// Package deadman implements a dead-man's switch for the portwatch monitor.
//
// If the scanner stops producing results — due to a crash, a hung goroutine,
// or a lost /proc mount — the Switch fires an AlertFunc after a configurable
// window so that operators are notified of the silence.
//
// Typical usage:
//
//	dm := deadman.New(30*time.Second, func(missed time.Duration) {
//		log.Printf("[WARN] no scan results for %s", missed)
//	})
//	dm.Start(ctx)
//
//	// Inside the scan loop:
//	dm.Reset()
package deadman
