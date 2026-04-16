// Package sigwatch bridges OS-level signals into the portwatch event bus.
//
// It translates SIGINT and SIGTERM into a "shutdown" event and SIGHUP into
// a "reload" event, allowing the rest of the daemon to react without
// importing os/signal directly.
//
// Typical use:
//
//	bus := eventbus.New()
//	watcher := sigwatch.New(bus)
//	go watcher.Run(ctx)
package sigwatch
