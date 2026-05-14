package aggregator_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/aggregator"
	"github.com/user/logslice/internal/parser"
)

func makeEntries() []parser.Entry {
	base := time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC)
	return []parser.Entry{
		{Timestamp: base, Level: "INFO", Message: "server started"},
		{Timestamp: base.Add(1 * time.Minute), Level: "INFO", Message: "request received"},
		{Timestamp: base.Add(2 * time.Minute), Level: "ERROR", Message: "connection failed"},
		{Timestamp: base.Add(3 * time.Minute), Level: "WARN", Message: "high memory usage"},
		{Timestamp: base.Add(4 * time.Minute), Level: "ERROR", Message: "timeout occurred"},
	}
}

func TestByLevel_Counts(t *testing.T) {
	entries := makeEntries()
	result := aggregator.ByLevel(entries)

	if result["INFO"].Count != 2 {
		t.Errorf("expected INFO count 2, got %d", result["INFO"].Count)
	}
	if result["ERROR"].Count != 2 {
		t.Errorf("expected ERROR count 2, got %d", result["ERROR"].Count)
	}
	if result["WARN"].Count != 1 {
		t.Errorf("expected WARN count 1, got %d", result["WARN"].Count)
	}
}

func TestByLevel_TimeRange(t *testing.T) {
	entries := makeEntries()
	result := aggregator.ByLevel(entries)

	info := result["INFO"]
	if !info.Last.After(info.First) {
		t.Error("expected INFO Last to be after First")
	}
}

func TestByLevel_Empty(t *testing.T) {
	result := aggregator.ByLevel(nil)
	if len(result) != 0 {
		t.Errorf("expected empty result for nil input, got %d entries", len(result))
	}
}

func TestByWindow_Basic(t *testing.T) {
	entries := makeEntries()
	result := aggregator.ByWindow(entries, 2*time.Minute)

	if len(result) == 0 {
		t.Fatal("expected non-empty result from ByWindow")
	}

	total := 0
	for _, s := range result {
		total += s.Count
	}
	if total != len(entries) {
		t.Errorf("expected total count %d, got %d", len(entries), total)
	}
}

func TestByWindow_ZeroDuration(t *testing.T) {
	entries := makeEntries()
	result := aggregator.ByWindow(entries, 0)
	if result != nil {
		t.Error("expected nil result for zero window duration")
	}
}

func TestByWindow_Empty(t *testing.T) {
	result := aggregator.ByWindow(nil, time.Minute)
	if result != nil {
		t.Error("expected nil result for empty entries")
	}
}
