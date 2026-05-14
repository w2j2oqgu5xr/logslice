// Package pipeline ties together parsing, filtering, aggregation, and export
// into a single composable processing chain.
package pipeline

import (
	"bufio"
	"io"

	"github.com/user/logslice/internal/exporter"
	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

// Config holds all configuration needed to run a pipeline.
type Config struct {
	Input   io.Reader
	Output  io.Writer
	Format  string
	Options filter.Options
}

// Pipeline represents a configured log processing pipeline.
type Pipeline struct {
	cfg Config
	exp *exporter.Exporter
}

// New creates a new Pipeline from the given Config.
func New(cfg Config) (*Pipeline, error) {
	exp, err := exporter.New(cfg.Output, cfg.Format)
	if err != nil {
		return nil, err
	}
	return &Pipeline{cfg: cfg, exp: exp}, nil
}

// Run reads lines from the configured input, parses them, applies filters,
// and writes the results to the configured output in the chosen format.
func (p *Pipeline) Run() (int, error) {
	var entries []parser.Entry

	scanner := bufio.NewScanner(p.cfg.Input)
	for scanner.Scan() {
		line := scanner.Text()
		entry, ok := parser.ParseLine(line)
		if !ok {
			continue
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	filtered := filter.Filter(entries, p.cfg.Options)

	if err := p.exp.Export(filtered); err != nil {
		return 0, err
	}

	return len(filtered), nil
}
