// Package eventbus implements a lightweight in-process publish/subscribe bus
// used by portwatch to decouple port-change detection from downstream consumers
// such as the notifier, history recorder, and rate-limiter.
//
// # Usage
//
//	bus := eventbus.New()
//
//	unsub := bus.Subscribe(func(e eventbus.Event) {
//		fmt.Printf("%s port %d/%s\n", e.Type, e.Port, e.Proto)
//	})
//	defer unsub()
//
//	bus.Publish(eventbus.Event{
//		Type:  eventbus.EventNew,
//		Port:  8080,
//		Proto: "tcp",
//	})
//
// Subscribers are called synchronously in registration order. If concurrent
// fan-out is required, wrap the handler body in a goroutine.
package eventbus
