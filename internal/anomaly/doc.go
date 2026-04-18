// Package anomaly provides detection of unusual port activity by evaluating
// alert-to-seen ratios stored in a scorecard. A Detector classifies port keys
// as suspicious or anomalous when their ratio crosses configurable thresholds,
// allowing operators to identify ports that trigger alerts disproportionately
// often relative to how frequently they are observed.
package anomaly
