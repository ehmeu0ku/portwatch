// Package metricsink provides a lightweight in-process metric accumulator.
//
// Counters are identified by string names and are safe for concurrent use.
// A Snapshot can be taken at any time to read all current values without
// blocking writers.
//
// Typical usage:
//
//	sink := metricsink.New()
//	sink.Inc("scans.total")
//	sink.Add("alerts.sent", 3)
//	fmt.Println(sink.Get("scans.total")) // 1
package metricsink
