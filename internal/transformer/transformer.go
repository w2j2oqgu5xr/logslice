// Package transformer provides entry-level transformation utilities
// for mutating log entries in a pipeline, such as uppercasing levels,
// normalizing messages, or remapping field values.
package transformer

import (
	"strings"

	"github.com/user/logslice/internal/parser"
)

// TransformFunc is a function that mutates a single log entry.
type TransformFunc func(entry parser.Entry) parser.Entry

// Transformer applies a chain of TransformFuncs to log entries.
type Transformer struct {
	fns []TransformFunc
}

// New returns a Transformer that applies the given transform functions
// in order to each entry.
func New(fns ...TransformFunc) *Transformer {
	return &Transformer{fns: fns}
}

// Apply runs all registered transform functions over the provided entries
// and returns the transformed slice.
func (t *Transformer) Apply(entries []parser.Entry) []parser.Entry {
	result := make([]parser.Entry, len(entries))
	for i, e := range entries {
		for _, fn := range t.fns {
			e = fn(e)
		}
		result[i] = e
	}
	return result
}

// NormalizeMessage returns a TransformFunc that trims whitespace and
// collapses consecutive internal spaces in each entry's message.
func NormalizeMessage() TransformFunc {
	return func(e parser.Entry) parser.Entry {
		parts := strings.Fields(e.Message)
		e.Message = strings.Join(parts, " ")
		return e
	}
}

// UppercaseMessage returns a TransformFunc that converts the entry
// message to uppercase.
func UppercaseMessage() TransformFunc {
	return func(e parser.Entry) parser.Entry {
		e.Message = strings.ToUpper(e.Message)
		return e
	}
}

// AddField returns a TransformFunc that adds or overwrites a key in
// the entry's Fields map.
func AddField(key, value string) TransformFunc {
	return func(e parser.Entry) parser.Entry {
		if e.Fields == nil {
			e.Fields = make(map[string]string)
		}
		e.Fields[key] = value
		return e
	}
}
