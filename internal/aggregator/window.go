package aggregator

import (
	"sort"
	"time"

	"github.com/user/logslice/internal/parser"
)

// WindowBucket represents a time window and the count of log entries within it.
type WindowBucket struct {
	Start time.Time
	End   time.Time
	Count int
}

// ByWindow groups log entries into fixed-duration time windows and returns
// a slice of WindowBucket values sorted by window start time.
// A zero or negative duration returns an empty slice.
func ByWindow(entries []parser.Entry, window time.Duration) []WindowBucket {
	if window <= 0 || len(entries) == 0 {
		return []WindowBucket{}
	}

	buckets := make(map[time.Time]*WindowBucket)

	for _, e := range entries {
		slot := e.Timestamp.Truncate(window)
		if b, ok := buckets[slot]; ok {
			b.Count++
		} else {
			buckets[slot] = &WindowBucket{
				Start: slot,
				End:   slot.Add(window),
				Count: 1,
			}
		}
	}

	result := make([]WindowBucket, 0, len(buckets))
	for _, b := range buckets {
		result = append(result, *b)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Start.Before(result[j].Start)
	})

	return result
}
