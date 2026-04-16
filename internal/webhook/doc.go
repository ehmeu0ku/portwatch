// Package webhook provides an HTTP notifier that POSTs alert events
// to a configured endpoint as JSON.
//
// Usage:
//
//	wh := webhook.New(webhook.Config{
//		URL:     "https://example.com/hooks/portwatch",
//		Timeout: 3 * time.Second,
//		Secret:  "s3cr3t",
//	})
//	if err := wh.Send(ctx, event); err != nil {
//		log.Println(err)
//	}
//
// The secret, when set, is forwarded verbatim in the
// X-Portwatch-Secret request header so the receiving server can
// authenticate the call.
package webhook
