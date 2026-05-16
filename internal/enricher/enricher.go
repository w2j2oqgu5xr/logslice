// Package enricher provides functionality for attaching additional metadata
// fields to log entries, such as hostname, environment tags, or custom labels.
package enricher

import (
	"os"

	"github.com/logslice/logslice/internal/parser"
)

// Enricher attaches additional fields to log entries.
type Enricher struct {
	fields map[string]string
}

// New creates a new Enricher with the given static fields.
func New(fields map[string]string) *Enricher {
	copy := make(map[string]string, len(fields))
	for k, v := range fields {
		copy[k] = v
	}
	return &Enricher{fields: copy}
}

// NewWithHostname creates an Enricher that automatically includes the system hostname.
func NewWithHostname(extra map[string]string) *Enricher {
	fields := make(map[string]string, len(extra)+1)
	for k, v := range extra {
		fields[k] = v
	}
	if host, err := os.Hostname(); err == nil {
		fields["hostname"] = host
	}
	return &Enricher{fields: fields}
}

// Enrich attaches the enricher's fields to a single entry's metadata.
func (e *Enricher) Enrich(entry parser.Entry) parser.Entry {
	if entry.Fields == nil {
		entry.Fields = make(map[string]string, len(e.fields))
	}
	for k, v := range e.fields {
		if _, exists := entry.Fields[k]; !exists {
			entry.Fields[k] = v
		}
	}
	return entry
}

// Apply enriches a slice of entries, returning a new slice.
func (e *Enricher) Apply(entries []parser.Entry) []parser.Entry {
	result := make([]parser.Entry, len(entries))
	for i, entry := range entries {
		result[i] = e.Enrich(entry)
	}
	return result
}
