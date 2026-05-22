// Package truncator provides utilities for truncating log entry messages
// that exceed a configurable maximum length, with optional suffix indicators.
package truncator

import "github.com/logslice/logslice/internal/parser"

// Truncator trims log entry messages that exceed a maximum byte length.
type Truncator struct {
	maxLen int
	suffix string
}

// New returns a Truncator that truncates messages longer than maxLen bytes.
// If maxLen is zero or negative, no truncation is applied.
// The suffix (e.g. "...") is appended when a message is truncated.
func New(maxLen int, suffix string) *Truncator {
	return &Truncator{
		maxLen: maxLen,
		suffix: suffix,
	}
}

// NewDefault returns a Truncator with a 256-byte limit and "..." suffix.
func NewDefault() *Truncator {
	return New(256, "...")
}

// Truncate shortens the message of a single entry if it exceeds the limit.
// Returns a new entry with the truncated message; the original is not modified.
func (t *Truncator) Truncate(e parser.Entry) parser.Entry {
	if t.maxLen <= 0 {
		return e
	}
	msg := e.Message
	if len(msg) <= t.maxLen {
		return e
	}
	cutAt := t.maxLen - len(t.suffix)
	if cutAt < 0 {
		cutAt = 0
	}
	e.Message = msg[:cutAt] + t.suffix
	return e
}

// Apply applies Truncate to every entry in the slice and returns the results.
func (t *Truncator) Apply(entries []parser.Entry) []parser.Entry {
	out := make([]parser.Entry, len(entries))
	for i, e := range entries {
		out[i] = t.Truncate(e)
	}
	return out
}
