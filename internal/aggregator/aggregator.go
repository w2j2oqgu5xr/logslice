// Package aggregator provides functionality for aggregating parsed log entries
// into summary statistics grouped by log level, time window, or message pattern.
package aggregator

import (
	"time"

	"github.com/user/logslice/internal/parser"
)

// Summary holds aggregated statistics for a group of log entries.
type Summary struct {
	Level    string            `json:"level"`
	Count    int               `json:"count"`
	First    time.Time         `json:"first"`
	Last     time.Time         `json:"last"`
	Messages []string          `json:"messages,omitempty"`
	Meta     map[string]string `json:"meta,omitempty"`
}

// ByLevel groups log entries by their level and returns a map of level to Summary.
func ByLevel(entries []parser.Entry) map[string]Summary {
	result := make(map[string]Summary)

	for _, e := range entries {
		s, ok := result[e.Level]
		if !ok {
			s = Summary{
				Level: e.Level,
				First: e.Timestamp,
				Last:  e.Timestamp,
			}
		}

		s.Count++
		s.Messages = append(s.Messages, e.Message)

		if e.Timestamp.Before(s.First) {
			s.First = e.Timestamp
		}
		if e.Timestamp.After(s.Last) {
			s.Last = e.Timestamp
		}

		result[e.Level] = s
	}

	return result
}

// ByWindow groups log entries into fixed-duration time windows and returns
// a slice of Summaries ordered by window start time.
func ByWindow(entries []parser.Entry, window time.Duration) []Summary {
	if len(entries) == 0 || window <= 0 {
		return nil
	}

	buckets := make(map[int64]*Summary)

	for _, e := range entries {
		key := e.Timestamp.Truncate(window).Unix()
		s, ok := buckets[key]
		if !ok {
			s = &Summary{
				First: e.Timestamp.Truncate(window),
				Last:  e.Timestamp,
			}
			buckets[key] = s
		}
		s.Count++
		if e.Timestamp.After(s.Last) {
			s.Last = e.Timestamp
		}
	}

	result := make([]Summary, 0, len(buckets))
	for _, s := range buckets {
		result = append(result, *s)
	}
	return result
}
