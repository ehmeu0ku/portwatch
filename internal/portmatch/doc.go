// Package portmatch provides flexible matching of scanner.PortState values
// against a list of user-defined expressions.
//
// An Expr can match by:
//   - exact port number
//   - inclusive port range (port_min..port_max)
//   - protocol ("tcp" or "udp")
//   - tag label attached to the state
//
// Multiple fields within a single Expr are ANDed together; all non-zero
// fields must be satisfied for the expression to match.
//
// A Matcher wraps a slice of Expr values and exposes AnyMatch (OR semantics)
// and AllMatch (AND semantics) helpers.
//
// Expression lists can be loaded from a JSON file via LoadFile, making it
// straightforward to drive match behaviour from configuration without
// recompiling the binary.
package portmatch
