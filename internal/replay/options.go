package replay

// Options configures replay behaviour.
type Options struct {
	// SpeedMult controls replay speed relative to real time.
	// 0 means no delay; 1.0 is real-time; 2.0 is double speed.
	SpeedMult float64

	// MaxEntries caps the total number of entries emitted. 0 means unlimited.
	MaxEntries int
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() Options {
	return Options{
		SpeedMult:  1.0,
		MaxEntries: 0,
	}
}

// OptionsBuilder provides a fluent API for constructing Options.
type OptionsBuilder struct {
	opts Options
}

// NewOptionsBuilder returns a builder initialised with defaults.
func NewOptionsBuilder() *OptionsBuilder {
	return &OptionsBuilder{opts: DefaultOptions()}
}

// WithSpeedMult sets the speed multiplier.
func (b *OptionsBuilder) WithSpeedMult(m float64) *OptionsBuilder {
	b.opts.SpeedMult = m
	return b
}

// WithMaxEntries sets the maximum number of entries to emit.
func (b *OptionsBuilder) WithMaxEntries(n int) *OptionsBuilder {
	b.opts.MaxEntries = n
	return b
}

// Build returns the constructed Options.
func (b *OptionsBuilder) Build() Options {
	return b.opts
}
