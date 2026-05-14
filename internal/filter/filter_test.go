package filter_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

var baseTime = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

func makeEntries() []parser.Entry {
	return []parser.Entry{
		{Timestamp: baseTime, Level: "INFO", Message: "server started"},
		{Timestamp: baseTime.Add(time.Minute), Level: "WARN", Message: "high memory usage"},
		{Timestamp: baseTime.Add(2 * time.Minute), Level: "ERROR", Message: "connection refused"},
		{Timestamp: baseTime.Add(3 * time.Minute), Level: "INFO", Message: "request received"},
	}
}

func TestFilter_ByLevel(t *testing.T) {
	entries := makeEntries()
	result := filter.Filter(entries, filter.Options{Level: "ERROR"})
	if len(result) != 1 || result[0].Level != "ERROR" {
		t.Errorf("expected 1 ERROR entry, got %d", len(result))
	}
}

func TestFilter_ByTimeRange(t *testing.T) {
	entries := makeEntries()
	opts := filter.Options{
		StartTime: baseTime.Add(time.Minute),
		EndTime:   baseTime.Add(2 * time.Minute),
	}
	result := filter.Filter(entries, opts)
	if len(result) != 2 {
		t.Errorf("expected 2 entries in range, got %d", len(result))
	}
}

func TestFilter_ByMessage(t *testing.T) {
	entries := makeEntries()
	result := filter.Filter(entries, filter.Options{MsgContains: "connection"})
	if len(result) != 1 || result[0].Message != "connection refused" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestFilter_NoOptions(t *testing.T) {
	entries := makeEntries()
	result := filter.Filter(entries, filter.Options{})
	if len(result) != len(entries) {
		t.Errorf("expected all %d entries, got %d", len(entries), len(result))
	}
}

func TestFilter_CaseInsensitiveLevel(t *testing.T) {
	entries := makeEntries()
	result := filter.Filter(entries, filter.Options{Level: "warn"})
	if len(result) != 1 || result[0].Level != "WARN" {
		t.Errorf("expected 1 WARN entry (case-insensitive), got %d", len(result))
	}
}

func TestFilter_EmptyInput(t *testing.T) {
	result := filter.Filter(nil, filter.Options{Level: "INFO"})
	if result != nil && len(result) != 0 {
		t.Errorf("expected empty result for nil input, got %+v", result)
	}
}
