// Package snapshot captures and persists point-in-time views of active
// port states observed by portwatch.
//
// A Snapshot records the full list of PortState values seen at a given
// moment along with the UTC timestamp of the capture. Snapshots can be
// saved to and loaded from JSON files, making them useful for:
//
//   - Offline diffing between two observation windows.
//   - Archiving the port landscape before a deployment.
//   - Feeding historical data back into the monitor for replay.
//
// Usage:
//
//	snap := snapshot.New(states)
//	if err := snap.Save("/var/lib/portwatch/latest.json"); err != nil {
//		log.Fatal(err)
//	}
//
//	loaded, err := snapshot.Load("/var/lib/portwatch/latest.json")
package snapshot
