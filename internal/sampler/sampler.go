// Package sampler provides log entry sampling strategies for reducing
// high-volume log streams to a representative subset.
package sampler

import (
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/user/logslice/internal/parser"
)

// Strategy defines how entries are selected from a stream.
type Strategy int

const (
	// StrategyRandom samples entries with a given probability.
	StrategyRandom Strategy = iota
	// StrategyNth keeps every Nth entry.
	StrategyNth
)

// Options configures the sampler behaviour.
type Options struct {
	Strategy   Strategy
	Rate       float64 // used by StrategyRandom: 0.0–1.0
	N          int     // used by StrategyNth: keep every N-th entry
}

// Sampler filters a slice of log entries according to the configured strategy.
type Sampler struct {
	opts    Options
	counter atomic.Int64
	rng     *rand.Rand
}

// New returns a Sampler configured with opts.
// Invalid options (e.g. Rate outside [0,1] or N < 1) are clamped to safe defaults.
func New(opts Options) *Sampler {
	if opts.Rate < 0 {
		opts.Rate = 0
	}
	if opts.Rate > 1 {
		opts.Rate = 1
	}
	if opts.N < 1 {
		opts.N = 1
	}
	return &Sampler{
		opts: opts,
		rng:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Sample returns the subset of entries selected by the configured strategy.
func (s *Sampler) Sample(entries []parser.Entry) []parser.Entry {
	if len(entries) == 0 {
		return nil
	}
	switch s.opts.Strategy {
	case StrategyNth:
		return s.sampleNth(entries)
	default:
		return s.sampleRandom(entries)
	}
}

func (s *Sampler) sampleRandom(entries []parser.Entry) []parser.Entry {
	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		if s.rng.Float64() < s.opts.Rate {
			out = append(out, e)
		}
	}
	return out
}

func (s *Sampler) sampleNth(entries []parser.Entry) []parser.Entry {
	out := make([]parser.Entry, 0, len(entries)/s.opts.N+1)
	for _, e := range entries {
		idx := s.counter.Add(1)
		if idx%int64(s.opts.N) == 0 {
			out = append(out, e)
		}
	}
	return out
}
