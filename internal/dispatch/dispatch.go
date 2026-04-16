// Package dispatch routes correlated events to configured notifiers,
// applying pagerules, rollup, and circuit-breaker protection.
package dispatch

import (
	"context"
	"fmt"

	"github.com/user/portwatch/internal/circuitbreaker"
	"github.com/user/portwatch/internal/correlator"
	"github.com/user/portwatch/internal/formatter"
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/pagerule"
)

// Dispatcher sends events to a notifier after applying rules and formatting.
type Dispatcher struct {
	rules   *pagerule.Engine
	fmt     *formatter.Formatter
	notify  notify.Notifier
	cb      *circuitbreaker.Breaker
}

// Config holds Dispatcher construction parameters.
type Config struct {
	Rules  *pagerule.Engine
	Format *formatter.Formatter
	Notify notify.Notifier
	Breaker *circuitbreaker.Breaker
}

// New returns a Dispatcher configured with the provided options.
func New(cfg Config) *Dispatcher {
	return &Dispatcher{
		rules:  cfg.Rules,
		fmt:    cfg.Format,
		notify: cfg.Notify,
		cb:     cfg.Breaker,
	}
}

// Dispatch evaluates the event against page rules and, if it should page,
// formats and forwards it through the notifier.
func (d *Dispatcher) Dispatch(ctx context.Context, ev correlator.Event) error {
	if !d.rules.ShouldPage(ev) {
		return nil
	}
	if !d.cb.Allow() {
		return fmt.Errorf("dispatch: circuit open, dropping event for port %d", ev.State.Port)
	}
	msg := d.fmt.Format(ev)
	if err := d.notify.Notify(ctx, msg); err != nil {
		d.cb.RecordFailure()
		return fmt.Errorf("dispatch: notify: %w", err)
	}
	d.cb.RecordSuccess()
	return nil
}
