// Package runner wires together the pipeline components based on a Config
// and executes the full logslice processing flow.
package runner

import (
	"fmt"
	"io"

	"github.com/yourorg/logslice/internal/config"
	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/pipeline"
	"github.com/yourorg/logslice/internal/reader"
)

// Runner executes the log processing pipeline using a resolved Config.
type Runner struct {
	cfg config.Config
	out io.Writer
}

// New creates a Runner for the given Config and output writer.
func New(cfg config.Config, out io.Writer) *Runner {
	return &Runner{cfg: cfg, out: out}
}

// Run builds the source, filter options, and pipeline, then streams output.
func (r *Runner) Run() error {
	src, err := r.buildSource()
	if err != nil {
		return fmt.Errorf("runner: build source: %w", err)
	}

	opts := r.buildFilterOptions()

	p := pipeline.New(src, opts, r.cfg.Format)
	if err := p.Run(r.out); err != nil {
		return fmt.Errorf("runner: pipeline: %w", err)
	}
	return nil
}

// buildSource constructs a LineSource from the configured input files or stdin.
func (r *Runner) buildSource() (pipeline.LineSource, error) {
	if len(r.cfg.InputFiles) == 0 {
		return reader.NewStdinSource(), nil
	}
	if len(r.cfg.InputFiles) == 1 {
		return reader.NewFileSource(r.cfg.InputFiles[0]), nil
	}
	sources := make([]pipeline.LineSource, 0, len(r.cfg.InputFiles))
	for _, f := range r.cfg.InputFiles {
		sources = append(sources, reader.NewFileSource(f))
	}
	return reader.NewMultiSource(sources...), nil
}

// buildFilterOptions converts Config fields into filter.Options.
func (r *Runner) buildFilterOptions() filter.Options {
	b := filter.NewOptionsBuilder()
	if r.cfg.Level != "" {
		b.WithLevel(r.cfg.Level)
	}
	if r.cfg.Message != "" {
		b.WithMessage(r.cfg.Message)
	}
	if !r.cfg.Since.IsZero() || !r.cfg.Until.IsZero() {
		b.WithTimeRange(r.cfg.Since, r.cfg.Until)
	}
	return b.Build()
}
