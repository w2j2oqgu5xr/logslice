// Package filter provides filtering capabilities for log entries.
package filter

import (
	"strings"
	"time"

	"github.com/user/logslice/internal/parser"
)

// Options holds the criteria used to filter log entries.
type Options struct {
	Level     string
	StartTime time.Time
	EndTime   time.Time
	MsgContains string
}

// Filter applies the given Options to a slice of log entries and returns
// only those entries that match all specified criteria.
func Filter(entries []parser.Entry, opts Options) []parser.Entry {
	var result []parser.Entry
	for _, e := range entries {
		if !matchLevel(e, opts.Level) {
			continue
		}
		if !matchTimeRange(e, opts.StartTime, opts.EndTime) {
			continue
		}
		if !matchMessage(e, opts.MsgContains) {
			continue
		}
		result = append(result, e)
	}
	return result
}

func matchLevel(e parser.Entry, level string) bool {
	if level == "" {
		return true
	}
	return strings.EqualFold(e.Level, level)
}

func matchTimeRange(e parser.Entry, start, end time.Time) bool {
	if !start.IsZero() && e.Timestamp.Before(start) {
		return false
	}
	if !end.IsZero() && e.Timestamp.After(end) {
		return false
	}
	return true
}

func matchMessage(e parser.Entry, substr string) bool {
	if substr == "" {
		return true
	}
	return strings.Contains(strings.ToLower(e.Message), strings.ToLower(substr))
}
