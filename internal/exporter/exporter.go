// Package exporter provides functionality for exporting parsed log entries
// to structured formats such as JSON and CSV.
package exporter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"

	"github.com/user/logslice/internal/parser"
)

// Format represents the output format for exported log entries.
type Format string

const (
	// FormatJSON exports entries as a JSON array.
	FormatJSON Format = "json"
	// FormatCSV exports entries as comma-separated values.
	FormatCSV Format = "csv"
)

// Exporter writes log entries to an io.Writer in a specified format.
type Exporter struct {
	writer io.Writer
	format Format
}

// New creates a new Exporter that writes to w in the given format.
func New(w io.Writer, format Format) *Exporter {
	return &Exporter{writer: w, format: format}
}

// Export writes the provided log entries to the underlying writer.
// It returns an error if the export fails.
func (e *Exporter) Export(entries []parser.Entry) error {
	switch e.format {
	case FormatJSON:
		return e.exportJSON(entries)
	case FormatCSV:
		return e.exportCSV(entries)
	default:
		return fmt.Errorf("unsupported format: %q", e.format)
	}
}

func (e *Exporter) exportJSON(entries []parser.Entry) error {
	enc := json.NewEncoder(e.writer)
	enc.SetIndent("", "  ")
	if err := enc.Encode(entries); err != nil {
		return fmt.Errorf("json export: %w", err)
	}
	return nil
}

func (e *Exporter) exportCSV(entries []parser.Entry) error {
	w := csv.NewWriter(e.writer)
	// Write header row.
	if err := w.Write([]string{"timestamp", "level", "message"}); err != nil {
		return fmt.Errorf("csv header: %w", err)
	}
	for _, entry := range entries {
		row := []string{
			entry.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
			entry.Level,
			entry.Message,
		}
		if err := w.Write(row); err != nil {
			return fmt.Errorf("csv row: %w", err)
		}
	}
	w.Flush()
	return w.Error()
}
