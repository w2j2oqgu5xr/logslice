package filter

import "time"

// OptionsBuilder provides a fluent API for constructing filter Options.
type OptionsBuilder struct {
	opts Options
}

// NewOptionsBuilder returns a new OptionsBuilder with default (zero) values.
func NewOptionsBuilder() *OptionsBuilder {
	return &OptionsBuilder{}
}

// WithLevel sets the log level filter.
func (b *OptionsBuilder) WithLevel(level string) *OptionsBuilder {
	b.opts.Level = level
	return b
}

// WithStartTime sets the earliest timestamp to include.
func (b *OptionsBuilder) WithStartTime(t time.Time) *OptionsBuilder {
	b.opts.StartTime = t
	return b
}

// WithEndTime sets the latest timestamp to include.
func (b *OptionsBuilder) WithEndTime(t time.Time) *OptionsBuilder {
	b.opts.EndTime = t
	return b
}

// WithMsgContains sets the message substring filter.
func (b *OptionsBuilder) WithMsgContains(substr string) *OptionsBuilder {
	b.opts.MsgContains = substr
	return b
}

// Build returns the constructed Options value.
func (b *OptionsBuilder) Build() Options {
	return b.opts
}
