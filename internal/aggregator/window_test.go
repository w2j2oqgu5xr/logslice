package aggregator

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeWindowEntries() []parser.Entry {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	return []parser.Entry{
		{Timestamp: base, Level: "INFO", Message: "a"},
		{Timestamp: base.Add(30 * time.Second), Level: "WARN", Message: "b"},
		{Timestamp: base.Add(90 * time.Second), Level: "INFO", Message: "c"},
		{Timestamp: base.Add(2*time.Minute + 10*time.Second), Level: "ERROR", Message: "d"},
		{Timestamp: base.Add(2*time.Minute + 50*time.Second), Level: "INFO", Message: "e"},
	}
}

func TestByWindow_Basic(t *testing.T) {
	entries := makeWindowEntries()
	results := ByWindow(entries, time.Minute)

	if len(results) == 0 {
		t.Fatal("expected non-empty results")
	}

	// First window should have INFO and WARN.
	first := results[0]
	if first.Counts["INFO"] != 1 {
		t.Errorf("expected 1 INFO in first window, got %d", first.Counts["INFO"])
	}
	if first.Counts["WARN"] != 1 {
		t.Errorf("expected 1 WARN in first window, got %d", first.Counts["WARN"])
	}
	if first.Total != 2 {
		t.Errorf("expected total 2 in first window, got %d", first.Total)
	}
}

func TestByWindow_WindowBoundary(t *testing.T) {
	entries := makeWindowEntries()
	results := ByWindow(entries, time.Minute)

	for _, r := range results {
		if !r.WindowEnd.Equal(r.WindowStart.Add(time.Minute)) {
			t.Errorf("window duration mismatch: start=%v end=%v", r.WindowStart, r.WindowEnd)
		}
	}
}

func TestByWindow_EmptyEntries(t *testing.T) {
	results := ByWindow(nil, time.Minute)
	if results != nil {
		t.Errorf("expected nil result for empty entries, got %v", results)
	}
}

func TestByWindow_ZeroDuration(t *testing.T) {
	entries := makeWindowEntries()
	results := ByWindow(entries, 0)
	if results != nil {
		t.Errorf("expected nil result for zero duration, got %v", results)
	}
}

func TestByWindow_SingleEntry(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	entries := []parser.Entry{
		{Timestamp: base, Level: "DEBUG", Message: "only one"},
	}
	results := ByWindow(entries, 5*time.Minute)
	if len(results) != 1 {
		t.Fatalf("expected 1 window, got %d", len(results))
	}
	if results[0].Counts["DEBUG"] != 1 {
		t.Errorf("expected 1 DEBUG, got %d", results[0].Counts["DEBUG"])
	}
}
