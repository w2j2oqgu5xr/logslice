package formatter

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/yourorg/logslice/internal/aggregator"
)

// SummaryFormatter renders aggregator.LevelSummary results as a table.
type SummaryFormatter struct {
	ShowPercentage bool
}

// DefaultSummaryFormatter returns a SummaryFormatter with defaults.
func DefaultSummaryFormatter() *SummaryFormatter {
	return &SummaryFormatter{ShowPercentage: true}
}

// FormatSummary writes a level summary table to w.
func (f *SummaryFormatter) FormatSummary(w io.Writer, summaries []aggregator.LevelSummary) error {
	if len(summaries) == 0 {
		_, err := fmt.Fprintln(w, "(no entries)")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)

	var total int
	for _, s := range summaries {
		total += s.Count
	}

	_, _ = fmt.Fprintln(tw, "LEVEL\tCOUNT\tFIRST\tLAST\tPCT")
	_, _ = fmt.Fprintln(tw, "-----\t-----\t-----\t----\t---")

	for _, s := range summaries {
		pct := 0.0
		if total > 0 {
			pct = float64(s.Count) / float64(total) * 100
		}

		line := fmt.Sprintf("%s\t%d\t%s\t%s",
			s.Level,
			s.Count,
			s.First.Format("15:04:05"),
			s.Last.Format("15:04:05"),
		)

		if f.ShowPercentage {
			line += fmt.Sprintf("\t%.1f%%", pct)
		} else {
			line += "\t-"
		}

		_, _ = fmt.Fprintln(tw, line)
	}

	return tw.Flush()
}
