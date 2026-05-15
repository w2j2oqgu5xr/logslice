// Package redactor provides log entry field redaction for sensitive data.
package redactor

import (
	"regexp"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// Rule defines a single redaction rule applied to log message text.
type Rule struct {
	Pattern     *regexp.Regexp
	Replacement string
}

// Redactor applies a set of rules to scrub sensitive data from log entries.
type Redactor struct {
	rules []Rule
}

// New creates a Redactor with the provided rules.
func New(rules []Rule) *Redactor {
	return &Redactor{rules: rules}
}

// NewDefault returns a Redactor pre-loaded with common sensitive-data patterns.
func NewDefault() *Redactor {
	return New([]Rule{
		{
			Pattern:     regexp.MustCompile(`(?i)(password|passwd|secret|token|api[_-]?key)=[^\s&]+`),
			Replacement: "$1=[REDACTED]",
		},
		{
			Pattern:     regexp.MustCompile(`\b[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}\b`),
			Replacement: "[EMAIL]",
		},
		{
			Pattern:     regexp.MustCompile(`\b(?:\d[ -]?){13,16}\b`),
			Replacement: "[CARD]",
		},
	})
}

// Redact returns a copy of the entry with sensitive data removed from its message.
func (r *Redactor) Redact(e parser.Entry) parser.Entry {
	msg := e.Message
	for _, rule := range r.rules {
		msg = rule.Pattern.ReplaceAllString(msg, rule.Replacement)
	}
	e.Message = strings.TrimSpace(msg)
	return e
}

// Apply redacts all entries in the slice, returning a new slice.
func (r *Redactor) Apply(entries []parser.Entry) []parser.Entry {
	out := make([]parser.Entry, len(entries))
	for i, e := range entries {
		out[i] = r.Redact(e)
	}
	return out
}
