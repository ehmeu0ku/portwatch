// Package pagerule provides a lightweight rule engine that maps
// correlator.Event values to paging actions.
//
// Rules are evaluated in order; the first matching rule wins. If no rule
// matches the event is suppressed by default, so operators only receive
// notifications for ports they explicitly care about.
//
// Example
//
//	rules := []pagerule.Rule{
//		{
//			MinSeverity: severity.Critical,
//			Action:      pagerule.ActionPage,
//		},
//		{
//			MinSeverity: severity.Info,
//			Tags:        []string{"http"},
//			Action:      pagerule.ActionPage,
//		},
//	}
//	ev := pagerule.New(rules)
//	action := ev.Evaluate(event)
package pagerule
