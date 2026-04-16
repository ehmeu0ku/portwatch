// Package dispatch wires together the pagerule engine, event formatter,
// circuit breaker, and notifier into a single Dispatch call.
//
// Typical usage:
//
//	d := dispatch.New(dispatch.Config{
//		Rules:   rules,
//		Format:  formatter.New(formatter.Text),
//		Notify:  notifier,
//		Breaker: circuitbreaker.New(circuitbreaker.DefaultConfig()),
//	})
//
//	for ev := range events {
//		if err := d.Dispatch(ctx, ev); err != nil {
//			log.Println(err)
//		}
//	}
package dispatch
