package filter_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/filter"
)

func TestOptionsBuilder_Defaults(t *testing.T) {
	opts := filter.NewOptionsBuilder().Build()
	if opts.Level != "" {
		t.Errorf("expected empty Level, got %q", opts.Level)
	}
	if !opts.StartTime.IsZero() {
		t.Errorf("expected zero StartTime")
	}
	if !opts.EndTime.IsZero() {
		t.Errorf("expected zero EndTime")
	}
	if opts.MsgContains != "" {
		t.Errorf("expected empty MsgContains, got %q", opts.MsgContains)
	}
}

func TestOptionsBuilder_Chaining(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	opts := filter.NewOptionsBuilder().
		WithLevel("WARN").
		WithStartTime(start).
		WithEndTime(end).
		WithMsgContains("disk").
		Build()

	if opts.Level != "WARN" {
		t.Errorf("expected Level WARN, got %q", opts.Level)
	}
	if !opts.StartTime.Equal(start) {
		t.Errorf("expected StartTime %v, got %v", start, opts.StartTime)
	}
	if !opts.EndTime.Equal(end) {
		t.Errorf("expected EndTime %v, got %v", end, opts.EndTime)
	}
	if opts.MsgContains != "disk" {
		t.Errorf("expected MsgContains 'disk', got %q", opts.MsgContains)
	}
}

func TestOptionsBuilder_IntegrationWithFilter(t *testing.T) {
	entries := makeEntries()
	opts := filter.NewOptionsBuilder().
		WithLevel("INFO").
		WithMsgContains("started").
		Build()

	result := filter.Filter(entries, opts)
	if len(result) != 1 {
		t.Errorf("expected 1 matching entry, got %d", len(result))
	}
}
