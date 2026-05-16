package parser

import "time"

// Entry represents a single parsed log line with structured fields.
type Entry struct {
	// Timestamp is the parsed time of the log event.
	Timestamp time.Time

	// Level is the severity level of the log entry.
	Level Level

	// Message is the human-readable log message.
	Message string

	// Raw is the original unparsed log line.
	Raw string

	// Fields holds optional key-value metadata attached to the entry,
	// either parsed from structured log output or added by enrichers.
	Fields map[string]string
}

// IsValid reports whether the entry has a non-zero timestamp and a valid level.
func (e Entry) IsValid() bool {
	return !e.Timestamp.IsZero() && e.Level.IsValid()
}
