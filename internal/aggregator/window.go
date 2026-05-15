package aggregator

import (
	"fmt"
	"time"

	"github.com/user/logslice/internal/parser"
)

// WindowResult holds aggregated log counts for a specific time window.
type WindowResult struct {
	WindowStart time.Time
	WindowEnd   time.Time
	Counts      map[string]int
	Total       int
}

// String returns a human-readable representation of a WindowResult.
func (w WindowResult) String() string {
	return fmt.Sprintf("[%s - %s] total=%d counts=%v",
		w.WindowStart.Format(time.RFC3339),
		w.WindowEnd.Format(time.RFC3339),
		w.Total,
		w.Counts,
	)
}

// ByWindow groups log entries into fixed-duration time windows and counts
// entries per log level within each window.
func ByWindow(entries []parser.Entry, window time.Duration) []WindowResult {
	if len(entries) == 0 || window <= 0 {
		return nil
	}

	// Determine the earliest timestamp to anchor windows.
	start := entries[0].Timestamp
	for _, e := range entries[1:] {
		if e.Timestamp.Before(start) {
			start = e.Timestamp
		}
	}

	// Truncate start to the window boundary.
	start = start.Truncate(window)

	// Bucket entries into windows.
	buckets := make(map[int64]*WindowResult)
	for _, e := range entries {
		offset := int64(e.Timestamp.Sub(start) / window)
		key := offset
		if _, ok := buckets[key]; !ok {
			wStart := start.Add(time.Duration(offset) * window)
			buckets[key] = &WindowResult{
				WindowStart: wStart,
				WindowEnd:   wStart.Add(window),
				Counts:      make(map[string]int),
			}
		}
		buckets[key].Counts[e.Level]++
		buckets[key].Total++
	}

	// Sort results by window start time.
	var results []WindowResult
	for i := int64(0); ; i++ {
		r, ok := buckets[i]
		if !ok {
			if i > int64(len(buckets)+10) {
				break
			}
			continue
		}
		results = append(results, *r)
	}
	return results
}
