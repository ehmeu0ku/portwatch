// Package portcategory assigns broad functional categories to port numbers.
//
// Built-in mappings cover common well-known ports. Callers may supply a
// custom map at construction time to override or extend the defaults.
//
// Example:
//
//	cl := portcategory.New(nil)
//	cat := cl.Classify(443) // → "web"
package portcategory
