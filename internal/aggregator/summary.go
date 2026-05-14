package aggregator

import (
	"sort"
	"time"

	"github.com/example/logslice/internal/parser"
)

// LevelSummary holds aggregated statistics for a single log level.
type LevelSummary struct {
	Level   string    `json:"level"   csv:"level"`
	Count   int       `json:"count"   csv:"count"`
	First   time.Time `json:"first"   csv:"first"`
	Last    time.Time `json:"last"    csv:"last"`
}

// SummarizeByLevel returns a slice of LevelSummary values, one per
// distinct log level found in entries, sorted alphabetically by level.
// If entries is empty the returned slice is non-nil but has length zero.
func SummarizeByLevel(entries []parser.Entry) []LevelSummary {
	type bucket struct {
		count int
		first time.Time
		last  time.Time
	}

	buckets := make(map[string]*bucket)

	for _, e := range entries {
		b, ok := buckets[e.Level]
		if !ok {
			buckets[e.Level] = &bucket{
				count: 1,
				first: e.Timestamp,
				last:  e.Timestamp,
			}
			continue
		}
		b.count++
		if e.Timestamp.Before(b.first) {
			b.first = e.Timestamp
		}
		if e.Timestamp.After(b.last) {
			b.last = e.Timestamp
		}
	}

	levels := make([]string, 0, len(buckets))
	for l := range buckets {
		levels = append(levels, l)
	}
	sort.Strings(levels)

	summaries := make([]LevelSummary, 0, len(levels))
	for _, l := range levels {
		b := buckets[l]
		summaries = append(summaries, LevelSummary{
			Level: l,
			Count: b.count,
			First: b.first,
			Last:  b.last,
		})
	}
	return summaries
}
