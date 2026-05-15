package formatter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/aggregator"
	"github.com/user/logslice/internal/formatter"
)

func makeWindowBuckets() []aggregator.WindowBucket {
	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	return []aggregator.WindowBucket{
		{Start: now, End: now.Add(5 * time.Minute), Count: 12},
		{Start: now.Add(5 * time.Minute), End: now.Add(10 * time.Minute), Count: 7},
		{Start: now.Add(10 * time.Minute), End: now.Add(15 * time.Minute), Count: 3},
	}
}

func TestFormatWindow_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	err := formatter.DefaultWindowFormatter(makeWindowBuckets(), &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Window Start") {
		t.Errorf("expected header 'Window Start' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Count") {
		t.Errorf("expected header 'Count' in output, got:\n%s", out)
	}
}

func TestFormatWindow_ContainsCounts(t *testing.T) {
	var buf bytes.Buffer
	_ = formatter.DefaultWindowFormatter(makeWindowBuckets(), &buf)
	out := buf.String()
	for _, expected := range []string{"12", "7", "3"} {
		if !strings.Contains(out, expected) {
			t.Errorf("expected count %q in output, got:\n%s", expected, out)
		}
	}
}

func TestFormatWindow_ContainsTimestamps(t *testing.T) {
	var buf bytes.Buffer
	_ = formatter.DefaultWindowFormatter(makeWindowBuckets(), &buf)
	out := buf.String()
	if !strings.Contains(out, "2024-01-15T10:00:00Z") {
		t.Errorf("expected RFC3339 timestamp in output, got:\n%s", out)
	}
}

func TestFormatWindow_EmptyBuckets(t *testing.T) {
	var buf bytes.Buffer
	err := formatter.DefaultWindowFormatter([]aggregator.WindowBucket{}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "No window data available.") {
		t.Errorf("expected empty message, got:\n%s", out)
	}
}

func TestFormatWindow_MultipleLines(t *testing.T) {
	var buf bytes.Buffer
	_ = formatter.DefaultWindowFormatter(makeWindowBuckets(), &buf)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + separator + 3 data rows = 5 lines
	if len(lines) != 5 {
		t.Errorf("expected 5 lines, got %d:\n%s", len(lines), buf.String())
	}
}
