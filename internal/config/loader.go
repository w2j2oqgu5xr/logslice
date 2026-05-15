package config

import (
	"flag"
	"time"
)

// LoadFromFlags parses command-line flags into a Config and validates it.
// It uses the provided FlagSet so callers can inject a custom one for testing.
func LoadFromFlags(fs *flag.FlagSet, args []string) (*Config, error) {
	var (
		level       string
		message     string
		startStr    string
		endStr      string
		format      string
		output      string
		aggregateBy string
		windowSec   int
	)

	fs.StringVar(&level, "level", "", "filter by log level (e.g. INFO, ERROR)")
	fs.StringVar(&message, "message", "", "filter by message substring")
	fs.StringVar(&startStr, "start", "", "start time filter (RFC3339)")
	fs.StringVar(&endStr, "end", "", "end time filter (RFC3339)")
	fs.StringVar(&format, "format", "json", "output format: json or csv")
	fs.StringVar(&output, "output", "", "output file path (default stdout)")
	fs.StringVar(&aggregateBy, "aggregate", "", "aggregate by: level or window")
	fs.IntVar(&windowSec, "window", 0, "window size in seconds (used with -aggregate=window)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	cfg := &Config{
		InputFiles:   fs.Args(),
		Level:        level,
		MessageQ:     message,
		OutputFormat: Format(format),
		OutputFile:   output,
		AggregateBy:  aggregateBy,
		WindowSize:   time.Duration(windowSec) * time.Second,
	}

	if startStr != "" {
		t, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			return nil, err
		}
		cfg.StartTime = t
	}

	if endStr != "" {
		t, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			return nil, err
		}
		cfg.EndTime = t
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
