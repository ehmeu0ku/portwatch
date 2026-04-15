package notify

import "fmt"

// MultiNotifier fans a single Message out to multiple Notifier backends.
// All backends are attempted; errors are collected and returned together.
type MultiNotifier struct {
	backends []Notifier
}

// NewMultiNotifier returns a MultiNotifier wrapping the provided backends.
func NewMultiNotifier(backends ...Notifier) *MultiNotifier {
	return &MultiNotifier{backends: backends}
}

// Name returns a combined name of all registered backends.
func (m *MultiNotifier) Name() string { return "multi" }

// Add appends a backend to the MultiNotifier at runtime.
func (m *MultiNotifier) Add(n Notifier) {
	m.backends = append(m.backends, n)
}

// Send delivers msg to every registered backend.
// If one or more backends fail, a combined error is returned.
func (m *MultiNotifier) Send(msg Message) error {
	var errs []error
	for _, b := range m.backends {
		if err := b.Send(msg); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", b.Name(), err))
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("notify: %d backend(s) failed: %v", len(errs), errs)
}
