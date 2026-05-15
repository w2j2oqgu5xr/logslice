// Package config provides configuration loading and validation for logslice.
package config

import (
	"errors"
	"strings"
	"time"
)

// Format represents the output format for exported data.
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// Config holds all runtime configuration for a logslice run.
type Config struct {
	// Input sources (file paths; empty means stdin)
	InputFiles []string

	// Filter options
	Level      string
	MessageQ   string
	StartTime  time.Time
	EndTime    time.Time

	// Output
	OutputFormat Format
	OutputFile   string // empty means stdout

	// Aggregation
	AggregateBy  string        // "level", "window", or ""
	WindowSize   time.Duration
}

// Validate checks that the Config fields are consistent and valid.
func (c *Config) Validate() error {
	switch c.OutputFormat {
	case FormatJSON, FormatCSV:
		// valid
	case "":
		c.OutputFormat = FormatJSON
	default:
		return errors.New("unsupported output format: " + string(c.OutputFormat))
	}

	agg := strings.ToLower(c.AggregateBy)
	if agg != "" && agg != "level" && agg != "window" {
		return errors.New("unsupported aggregate-by value: " + c.AggregateBy)
	}
	c.AggregateBy = agg

	if agg == "window" && c.WindowSize <= 0 {
		return errors.New("window aggregation requires a positive window-size")
	}

	if !c.StartTime.IsZero() && !c.EndTime.IsZero() && c.EndTime.Before(c.StartTime) {
		return errors.New("end-time must not be before start-time")
	}

	return nil
}
