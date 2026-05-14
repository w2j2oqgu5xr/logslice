// Package aggregator provides functions for summarizing and grouping
// parsed log entries into structured aggregate results.
//
// # Overview
//
// The aggregator package exposes two primary aggregation strategies:
//
//   - ByLevel: groups log entries by their severity level and returns
//     counts, optionally filtered by a time range.
//
//   - ByWindow: groups log entries into fixed-duration time windows,
//     counting entries that fall within each window boundary.
//
// # Usage
//
//	counts := aggregator.ByLevel(entries, start, end)
//	windows := aggregator.ByWindow(entries, time.Minute*5)
//
// Both functions accept a slice of parser.Entry values and return
// map-based summaries suitable for downstream export or display.
package aggregator
