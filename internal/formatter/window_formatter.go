package formatter

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/user/logslice/internal/aggregator"
)

// DefaultWindowFormatter formats a slice of WindowBucket results into a
// human-readable table-style text output.
func DefaultWindowFormatter(buckets []aggregator.WindowBucket, w io.Writer) error {
	if len(buckets) == 0 {
		_, err := fmt.Fprintln(w, "No window data available.")
		return err
	}

	header := fmt.Sprintf("%-30s %-30s %-10s", "Window Start", "Window End", "Count")
	separator := strings.Repeat("-", 72)

	if _, err := fmt.Fprintln(w, header); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, separator); err != nil {
		return err
	}

	for _, b := range buckets {
		line := fmt.Sprintf("%-30s %-30s %-10d",
			formatWindowTime(b.Start),
			formatWindowTime(b.End),
			b.Count,
		)
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}

	return nil
}

func formatWindowTime(t time.Time) string {
	if t.IsZero() {
		return "N/A"
	}
	return t.UTC().Format(time.RFC3339)
}
