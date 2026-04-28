// Package portfence provides a protocol-aware port access fence.
//
// A Fence holds an explicit allow-list of ports per protocol (tcp/udp).
// Any port that is not in the allow-list for its protocol is considered a
// violation. Protocols with no fence defined are left unrestricted, so the
// fence can be applied incrementally.
//
// Typical usage:
//
//	f := portfence.New()
//	f.Allow("tcp", 80)
//	f.Allow("tcp", 443)
//
//	for _, v := range f.Violations(currentStates) {
//		log.Printf("unexpected port: %v", v)
//	}
package portfence
