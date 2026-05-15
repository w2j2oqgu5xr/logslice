package formatter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/aggregator"
	"github.com/yourorg/logslice/internal/formatter"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(ts time.Time, level parser.Level, msg string) parser.Entry {
	return parser.Entry{Timestamp: ts, Level: level, Message: msg}
}

func TestFormatEntry_ContainsTimestamp(t *testing.T) {
	f := formatter.DefaultEntryFormatter()
	var buf bytes.Buffer
	ts := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	err := f.FormatEntry(&buf, makeEntry(ts, parser.LevelInfo, "hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "2024-01-15") {
		t.Errorf("expected timestamp in output, got: %s", buf.String())
	}
}

func TestFormatEntry_ContainsLevel(t *testing.T) {
	f := formatter.DefaultEntryFormatter()
	var buf bytes.Buffer
	ts := time.Now()
	_ = f.FormatEntry(&buf, makeEntry(ts, parser.LevelError, "boom"))
	if !strings.Contains(buf.String(), "[ERROR]") {
		t.Errorf("expected [ERROR] in output, got: %s", buf.String())
	}
}

func TestFormatEntry_ContainsMessage(t *testing.T) {
	f := formatter.DefaultEntryFormatter()
	var buf bytes.Buffer
	ts := time.Now()
	_ = f.FormatEntry(&buf, makeEntry(ts, parser.LevelWarn, "disk full"))
	if !strings.Contains(buf.String(), "disk full") {
		t.Errorf("expected message in output, got: %s", buf.String())
	}
}

func TestFormatEntries_MultipleLines(t *testing.T) {
	f := formatter.DefaultEntryFormatter()
	var buf bytes.Buffer
	base := time.Date(2024, 3, 1, 9, 0, 0, 0, time.UTC)
	entries := []parser.Entry{
		makeEntry(base, parser.LevelInfo, "start"),
		makeEntry(base.Add(time.Second), parser.LevelError, "fail"),
	}
	if err := f.FormatEntries(&buf, entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestFormatSummary_Empty(t *testing.T) {
	sf := formatter.DefaultSummaryFormatter()
	var buf bytes.Buffer
	if err := sf.FormatSummary(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no entries") {
		t.Errorf("expected 'no entries' message, got: %s", buf.String())
	}
}

func TestFormatSummary_ShowsLevelAndCount(t *testing.T) {
	sf := formatter.DefaultSummaryFormatter()
	var buf bytes.Buffer
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	summaries := []aggregator.LevelSummary{
		{Level: "INFO", Count: 5, First: base, Last: base.Add(time.Minute)},
		{Level: "ERROR", Count: 2, First: base, Last: base.Add(time.Minute)},
	}
	if err := sf.FormatSummary(&buf, summaries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "INFO") || !strings.Contains(out, "ERROR") {
		t.Errorf("expected level names in output, got: %s", out)
	}
	if !strings.Contains(out, "71.4%") {
		t.Errorf("expected percentage in output, got: %s", out)
	}
}
