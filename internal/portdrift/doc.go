// Package portdrift detects address-binding drift for monitored ports.
//
// A "drift" occurs when a port that was previously bound to a specific
// address (e.g. 127.0.0.1) is later observed bound to a different address
// (e.g. 0.0.0.0). This can indicate a misconfiguration or a deliberate
// attempt to expose a previously internal service.
//
// Usage:
//
//	det := portdrift.New()
//	events := det.Observe(currentStates)
//	for _, e := range events {
//		log.Println(e)
//	}
package portdrift
