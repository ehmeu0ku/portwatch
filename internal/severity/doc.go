// Package severity provides a simple three-level classification scheme
// (Info, Warning, Critical) for port-listener events detected by portwatch.
//
// Classification rules:
//
//   - Ports within the privileged range (≤ privilegedMax, default 1023) are
//     always rated Critical, because binding them typically requires elevated
//     privileges.
//
//   - Ports above the privileged threshold that match a known service tag
//     (via the tagger package) are rated Info — they are expected listeners.
//
//   - Ports above the privileged threshold with no known tag are rated
//     Warning, indicating an unrecognised listener that may warrant attention.
//
// Usage:
//
//	classifier := severity.New(1023, myTagger)
//	level := classifier.Classify(port)
//	fmt.Println(level) // "INFO", "WARNING", or "CRITICAL"
package severity
