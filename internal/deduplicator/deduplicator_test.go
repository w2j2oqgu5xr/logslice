package deduplicator_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/deduplicator"
	"github.com/logslice/logslice/internal/parser"
)

func makeEntries(msgs []string, level parser.Level) []parser.Entry {
	var entries []parser.Entry
	for _, m := range msgs {
		entries = append(entries, parser.Entry{
			Timestamp: time.Now(),
			Level:     level,
			Message:   m,
		})
	}
	return entries
}

func TestDeduplicator_RemovesDuplicates(t *testing.T) {
	d := deduplicator.New(1 * time.Minute)
	entries := makeEntries([]string{"disk full", "disk full", "cpu high"}, parser.LevelWarn)

	result := d.Filter(entries)
	if len(result) != 2 {
		t.Fatalf("expected 2 unique entries, got %d", len(result))
	}
}

func TestDeduplicator_AllowsDifferentLevels(t *testing.T) {
	d := deduplicator.New(1 * time.Minute)
	entries := []parser.Entry{
		{Timestamp: time.Now(), Level: parser.LevelInfo, Message: "started"},
		{Timestamp: time.Now(), Level: parser.LevelError, Message: "started"},
	}

	result := d.Filter(entries)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries (different levels), got %d", len(result))
	}
}

func TestDeduplicator_AllowsAfterWindowExpires(t *testing.T) {
	d := deduplicator.New(50 * time.Millisecond)

	first := makeEntries([]string{"timeout"}, parser.LevelError)
	d.Filter(first)

	time.Sleep(100 * time.Millisecond)

	second := makeEntries([]string{"timeout"}, parser.LevelError)
	result := d.Filter(second)
	if len(result) != 1 {
		t.Fatalf("expected entry to reappear after window, got %d", len(result))
	}
}

func TestDeduplicator_Reset(t *testing.T) {
	d := deduplicator.New(1 * time.Minute)
	entries := makeEntries([]string{"crash"}, parser.LevelError)
	d.Filter(entries)
	d.Reset()

	result := d.Filter(entries)
	if len(result) != 1 {
		t.Fatalf("expected entry after reset, got %d", len(result))
	}
}

func TestDeduplicator_EmptyInput(t *testing.T) {
	d := deduplicator.New(1 * time.Minute)
	result := d.Filter(nil)
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}

func TestDeduplicator_MultipleCallsAccumulateState(t *testing.T) {
	d := deduplicator.New(1 * time.Minute)

	// First call seeds the deduplicator with "retry"
	first := makeEntries([]string{"retry"}, parser.LevelWarn)
	d.Filter(first)

	// Second call with the same message should be filtered out
	second := makeEntries([]string{"retry", "new event"}, parser.LevelWarn)
	result := d.Filter(second)
	if len(result) != 1 {
		t.Fatalf("expected 1 entry (duplicate filtered across calls), got %d", len(result))
	}
	if result[0].Message != "new event" {
		t.Fatalf("expected 'new event', got %q", result[0].Message)
	}
}
