// Package porttrend provides a lightweight trend tracker for port presence
// across successive scans. It maintains a sliding window of boolean samples
// (seen / not-seen) per port key and derives a Rising, Falling, or Stable
// direction signal that downstream components can use to suppress noisy
// transient alerts or escalate persistently new listeners.
package porttrend
