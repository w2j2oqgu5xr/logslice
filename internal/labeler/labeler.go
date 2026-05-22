// Package labeler provides functionality for tagging log entries
// with user-defined labels based on matching rules.
package labeler

import (
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// Rule defines a matching condition and the label to apply.
type Rule struct {
	// Label is the tag value to assign when the rule matches.
	Label string
	// Level restricts matching to a specific log level. Empty means any level.
	Level parser.Level
	// MessageContains is a substring that must appear in the message.
	MessageContains string
}

// Labeler applies label rules to log entries, storing results in entry Fields.
type Labeler struct {
	rules []Rule
	field string
}

// New creates a Labeler that writes matched labels into the given field name.
func New(field string, rules []Rule) *Labeler {
	return &Labeler{field: field, rules: rules}
}

// NewDefault creates a Labeler that writes labels into the "label" field.
func NewDefault(rules []Rule) *Labeler {
	return New("label", rules)
}

// Label evaluates all rules against the entry and returns a copy with the
// first matching label applied. If no rule matches the entry is unchanged.
func (l *Labeler) Label(e parser.Entry) parser.Entry {
	for _, r := range l.rules {
		if r.Level != 0 && e.Level != r.Level {
			continue
		}
		if r.MessageContains != "" && !strings.Contains(e.Message, r.MessageContains) {
			continue
		}
		if e.Fields == nil {
			e.Fields = make(map[string]string)
		}
		e.Fields[l.field] = r.Label
		return e
	}
	return e
}

// Apply labels all entries in the slice, returning a new slice.
func (l *Labeler) Apply(entries []parser.Entry) []parser.Entry {
	out := make([]parser.Entry, len(entries))
	for i, e := range entries {
		out[i] = l.Label(e)
	}
	return out
}
