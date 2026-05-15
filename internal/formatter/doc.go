// Package formatter provides human-readable text rendering for log entries
// and aggregation summaries produced by logslice.
//
// # Overview
//
// The formatter package complements the exporter package (which handles
// structured JSON/CSV output) by offering plain-text, tabular representations
// suitable for terminal display.
//
// # Entry Formatting
//
// Use [DefaultEntryFormatter] to obtain a pre-configured renderer that prints
// each log entry as a timestamped, level-tagged line:
//
//	f := formatter.DefaultEntryFormatter()
//	f.FormatEntries(os.Stdout, entries)
//
// # Summary Formatting
//
// Use [DefaultSummaryFormatter] to render a [aggregator.LevelSummary] slice
// as an aligned table with counts and percentages:
//
//	sf := formatter.DefaultSummaryFormatter()
//	sf.FormatSummary(os.Stdout, summaries)
package formatter
