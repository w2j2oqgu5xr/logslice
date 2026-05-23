// Package masker provides field-level masking for log entries,
// replacing specific field values with a fixed placeholder string.
package masker

import (
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

const defaultPlaceholder = "[MASKED]"

// Masker replaces the values of named fields in log entry metadata
// with a configurable placeholder string.
type Masker struct {
	fields      map[string]struct{}
	placeholder string
}

// New returns a Masker that masks the given field names using the
// default placeholder.
func New(fields []string) *Masker {
	return NewWithPlaceholder(fields, defaultPlaceholder)
}

// NewWithPlaceholder returns a Masker that masks the given field names
// using the supplied placeholder string.
func NewWithPlaceholder(fields []string, placeholder string) *Masker {
	set := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		set[strings.ToLower(f)] = struct{}{}
	}
	return &Masker{fields: set, placeholder: placeholder}
}

// Mask returns a copy of the entry with targeted fields replaced by
// the placeholder. Fields not present in the entry are left unchanged.
func (m *Masker) Mask(e parser.Entry) parser.Entry {
	if len(e.Fields) == 0 || len(m.fields) == 0 {
		return e
	}
	masked := make(map[string]string, len(e.Fields))
	for k, v := range e.Fields {
		if _, ok := m.fields[strings.ToLower(k)]; ok {
			masked[k] = m.placeholder
		} else {
			masked[k] = v
		}
	}
	e.Fields = masked
	return e
}

// Apply masks every entry in the slice and returns the results.
func (m *Masker) Apply(entries []parser.Entry) []parser.Entry {
	out := make([]parser.Entry, len(entries))
	for i, e := range entries {
		out[i] = m.Mask(e)
	}
	return out
}
