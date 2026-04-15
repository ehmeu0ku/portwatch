// Package notify defines the Notifier interface and built-in backends used by
// portwatch to deliver alerts when unexpected port activity is detected.
//
// # Backends
//
// LogNotifier writes human-readable lines to any io.Writer (default: stdout).
// MultiNotifier fans a single message out to multiple backends simultaneously.
//
// # Usage
//
//	// Single backend
//	n := notify.NewLogNotifier(os.Stderr)
//	n.Send(notify.Message{
//		Level: notify.LevelAlert,
//		Title: "new listener",
//		Body:  "TCP :8080",
//		Timestamp: time.Now(),
//	})
//
//	// Fan-out
//	multi := notify.NewMultiNotifier(n1, n2)
//	multi.Send(msg)
package notify
