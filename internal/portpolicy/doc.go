// Package portpolicy provides a rule-based allow/deny engine for listening
// ports. Rules are evaluated in declaration order; the first matching rule
// wins. When no rule matches the policy's configured default action applies.
//
// Policies can be constructed programmatically via New/Add or loaded from a
// JSON file with LoadFile. The Enforcer adapter bridges portpolicy with the
// correlator event stream so that downstream components can gate alerts on
// policy decisions.
package portpolicy
