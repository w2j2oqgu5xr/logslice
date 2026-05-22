// Package splitter provides functionality for splitting a stream of log
// entries into multiple named buckets based on configurable field criteria.
package splitter

import (
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Rule defines a named split target and the condition that routes an entry
// into that bucket.
type Rule struct {
	// Name is the label assigned to the bucket.
	Name string
	// Level, when non-empty, matches entries whose level equals this value
	// (case-insensitive).
	Level string
	// MessageContains, when non-empty, matches entries whose message contains
	// this substring (case-insensitive).
	MessageContains string
}

// Splitter partitions log entries into named buckets according to a
// prioritised list of rules. Entries that match no rule are placed in the
// "default" bucket.
type Splitter struct {
	rules []Rule
}

// New returns a Splitter that evaluates rules in the order they are provided.
// The first matching rule wins.
func New(rules []Rule) *Splitter {
	return &Splitter{rules: rules}
}

// Split distributes entries across named buckets. Every entry is placed in
// exactly one bucket. Entries that match no rule land in "default".
func (s *Splitter) Split(entries []parser.Entry) map[string][]parser.Entry {
	buckets := make(map[string][]parser.Entry)
	for _, e := range entries {
		name := s.match(e)
		buckets[name] = append(buckets[name], e)
	}
	return buckets
}

// match returns the name of the first rule that matches e, or "default".
func (s *Splitter) match(e parser.Entry) string {
	for _, r := range s.rules {
		if r.Level != "" && !strings.EqualFold(e.Level.String(), r.Level) {
			continue
		}
		if r.MessageContains != "" &&
			!strings.Contains(strings.ToLower(e.Message), strings.ToLower(r.MessageContains)) {
			continue
		}
		return r.Name
	}
	return "default"
}
