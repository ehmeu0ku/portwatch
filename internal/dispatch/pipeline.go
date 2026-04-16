package dispatch

import (
	"context"
	"log"

	"github.com/user/portwatch/internal/correlator"
)

// RunPipeline reads events from ch and dispatches each one until ctx is done.
// Errors are logged but do not halt the pipeline.
func RunPipeline(ctx context.Context, d *Dispatcher, ch <-chan correlator.Event) {
	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-ch:
			if !ok {
				return
			}
			if err := d.Dispatch(ctx, ev); err != nil {
				log.Printf("dispatch pipeline: %v", err)
			}
		}
	}
}
