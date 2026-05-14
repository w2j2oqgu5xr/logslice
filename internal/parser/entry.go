package parser

import "time"

// LogLevel represents the severity level of a log entry.
type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
	LevelFatal LogLevel = "FATAL"
)

// Entry represents a single parsed log line with structured fields.
type Entry struct {
	Timestamp time.Time         `json:"timestamp" csv:"timestamp"`
	Level     LogLevel          `json:"level"     csv:"level"`
	Message   string            `json:"message"   csv:"message"`
	Source    string            `json:"source"    csv:"source"`
	Fields    map[string]string `json:"fields"    csv:"-"`
	Raw       string            `json:"-"         csv:"-"`
}

// IsValid returns true if the entry has a non-zero timestamp and non-empty message.
func (e *Entry) IsValid() bool {
	return !e.Timestamp.IsZero() && e.Message != ""
}

// HasField returns true if the given key exists in the entry's Fields map.
func (e *Entry) HasField(key string) bool {
	if e.Fields == nil {
		return false
	}
	_, ok := e.Fields[key]
	return ok
}
