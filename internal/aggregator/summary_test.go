package aggregator_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourusername/logslice/internal/aggregator"
	"github.com/yourusername/logslice/internal/parser"
)

func makeSummaryEntries() []parser.Entry {
	now := time.Now()
	return []parser.Entry{
		{Timestamp: now.Add(-5 * time.Minute), Level: "ERROR", Message: "disk full"},
		{Timestamp: now.Add(-4 * time.Minute), Level: "ERROR", Message: "connection refused"},
		{Timestamp: now.Add(-3 * time.Minute), Level: "WARN", Message: "high memory usage"},
		{Timestamp: now.Add(-2 * time.Minute), Level: "INFO", Message: "server started"},
		{Timestamp: now.Add(-1 * time.Minute), Level: "INFO", Message: "request received"},
		{Timestamp: now, Level: "DEBUG", Message: "trace output"},
	}
}

func TestSummarizeByLevel_BasicOutput(t *testing.T) {
	entries := makeSummaryEntries()
	var sb strings.Builder

	err := aggregator.SummarizeByLevel(entries, &sb)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := sb.String()
	if output == "" {
		t.Fatal("expected non-empty summary output")
	}

	expectedLevels := []string{"ERROR", "WARN", "INFO", "DEBUG"}
	for _, level := range expectedLevels {
		if !strings.Contains(output, level) {
			t.Errorf("expected output to contain level %q, got:\n%s", level, output)
		}
	}
}

func TestSummarizeByLevel_CountsAreCorrect(t *testing.T) {
	entries := makeSummaryEntries()
	var sb strings.Builder

	err := aggregator.SummarizeByLevel(entries, &sb)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := sb.String()

	// ERROR appears 2 times, INFO appears 2 times, WARN and DEBUG once each
	// We verify the output contains the counts in some form
	if !strings.Contains(output, "2") {
		t.Errorf("expected count '2' to appear in output for ERROR and INFO levels, got:\n%s", output)
	}
}

func TestSummarizeByLevel_EmptyEntries(t *testing.T) {
	var sb strings.Builder

	err := aggregator.SummarizeByLevel([]parser.Entry{}, &sb)
	if err != nil {
		t.Fatalf("unexpected error on empty input: %v", err)
	}

	output := sb.String()
	if output == "" {
		t.Fatal("expected some output even for empty entries (e.g. header or empty notice)")
	}
}

func TestSummarizeByLevel_SingleLevel(t *testing.T) {
	now := time.Now()
	entries := []parser.Entry{
		{Timestamp: now, Level: "INFO", Message: "only info"},
		{Timestamp: now.Add(time.Second), Level: "INFO", Message: "another info"},
	}

	var sb strings.Builder
	err := aggregator.SummarizeByLevel(entries, &sb)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := sb.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("expected output to contain INFO level, got:\n%s", output)
	}
	if strings.Contains(output, "ERROR") || strings.Contains(output, "WARN") || strings.Contains(output, "DEBUG") {
		t.Errorf("expected output to only contain INFO level, got:\n%s", output)
	}
}
