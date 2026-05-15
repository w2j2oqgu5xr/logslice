// Package filter provides log entry filtering based on configurable criteria.
package filter

import "time"

// Options holds the configuration for filtering log entries.
// Zero values mean "no filter applied" for that field.
type Options struct {
	// Level filters entries by exact log level (e.g., "ERROR", "INFO").
	// Empty string means all levels are accepted.
	Level string

	// From filters entries to those at or after this time.
	// Zero value means no lower bound.
	From time.Time

	// To filters entries to those at or before this time.
	// Zero value means no upper bound.
	To time.Time

	// MessageContains filters entries whose message contains this substring.
	// Empty string means all messages are accepted.
	MessageContains string
}

// IsZero reports whether no filter criteria have been set.
func (o Options) IsZero() bool {
	return o.Level == "" &&
		o.From.IsZero() &&
		o.To.IsZero() &&
		o.MessageContains == ""
}

// WithLevel returns a copy of Options with the Level field set.
func (o Options) WithLevel(level string) Options {
	o.Level = level
	return o
}

// WithTimeRange returns a copy of Options with From and To set.
func (o Options) WithTimeRange(from, to time.Time) Options {
	o.From = from
	o.To = to
	return o
}

// WithMessage returns a copy of Options with MessageContains set.
func (o Options) WithMessage(substr string) Options {
	o.MessageContains = substr
	return o
}
