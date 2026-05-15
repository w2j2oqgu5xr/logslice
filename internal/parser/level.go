package parser

import "strings"

// Level represents the severity level of a log entry.
type Level int

const (
	// LevelUnknown is used when the level cannot be determined.
	LevelUnknown Level = iota
	// LevelDebug represents debug-level log entries.
	LevelDebug
	// LevelInfo represents informational log entries.
	LevelInfo
	// LevelWarn represents warning-level log entries.
	LevelWarn
	// LevelError represents error-level log entries.
	LevelError
	// LevelFatal represents fatal-level log entries.
	LevelFatal
)

// levelNames maps string representations to Level constants.
var levelNames = map[string]Level{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"warning": LevelWarn,
	"error": LevelError,
	"fatal": LevelFatal,
}

// levelStrings maps Level constants to canonical string representations.
var levelStrings = map[Level]string{
	LevelUnknown: "UNKNOWN",
	LevelDebug:   "DEBUG",
	LevelInfo:    "INFO",
	LevelWarn:    "WARN",
	LevelError:   "ERROR",
	LevelFatal:   "FATAL",
}

// ParseLevel converts a string to a Level constant.
// Comparison is case-insensitive. Returns LevelUnknown if unrecognized.
func ParseLevel(s string) Level {
	if l, ok := levelNames[strings.ToLower(strings.TrimSpace(s))]; ok {
		return l
	}
	return LevelUnknown
}

// String returns the canonical string representation of a Level.
func (l Level) String() string {
	if s, ok := levelStrings[l]; ok {
		return s
	}
	return "UNKNOWN"
}

// IsValid reports whether the level is a recognised, non-unknown level.
func (l Level) IsValid() bool {
	return l != LevelUnknown
}
