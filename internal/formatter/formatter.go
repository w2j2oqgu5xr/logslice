// Package formatter provides utilities for rendering log entries
// and aggregation results into human-readable text output.
package formatter

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// EntryFormatter renders log entries as formatted text lines.
type EntryFormatter struct {
	TimestampFormat string
	ShowLevel       bool
	ShowMessage     bool
}

// DefaultEntryFormatter returns an EntryFormatter with sensible defaults.
func DefaultEntryFormatter() *EntryFormatter {
	return &EntryFormatter{
		TimestampFormat: time.RFC3339,
		ShowLevel:       true,
		ShowMessage:     true,
	}
}

// FormatEntry writes a single log entry as a formatted line to w.
func (f *EntryFormatter) FormatEntry(w io.Writer, e parser.Entry) error {
	parts := []string{}

	parts = append(parts, e.Timestamp.Format(f.TimestampFormat))

	if f.ShowLevel {
		parts = append(parts, fmt.Sprintf("[%s]", strings.ToUpper(e.Level.String())))
	}

	if f.ShowMessage {
		parts = append(parts, e.Message)
	}

	_, err := fmt.Fprintln(w, strings.Join(parts, " "))
	return err
}

// FormatEntries writes multiple log entries to w using the formatter.
func (f *EntryFormatter) FormatEntries(w io.Writer, entries []parser.Entry) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, e := range entries {
		if err := f.FormatEntry(tw, e); err != nil {
			return err
		}
	}
	return tw.Flush()
}
