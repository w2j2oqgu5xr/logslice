// Package config handles runtime configuration for logslice.
//
// It provides the Config struct which centralises all settings that control
// how logslice reads, filters, aggregates, and exports log data.
//
// Configuration can be populated from command-line flags via LoadFromFlags,
// which accepts a *flag.FlagSet so that callers (including tests) can supply
// their own flag set without touching os.Args.
//
// Supported flags:
//
//	-level      filter entries by log level (e.g. INFO, WARN, ERROR)
//	-message    filter entries whose message contains this substring
//	-start      include only entries at or after this RFC3339 timestamp
//	-end        include only entries at or before this RFC3339 timestamp
//	-format     output format: "json" (default) or "csv"
//	-output     write output to this file path instead of stdout
//	-aggregate  aggregate mode: "level" or "window"
//	-window     window size in seconds when -aggregate=window
//
// After parsing, Validate is called automatically to ensure the resulting
// Config is internally consistent before it is returned to the caller.
package config
