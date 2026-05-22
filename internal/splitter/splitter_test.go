package splitter_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/splitter"
)

func makeEntries() []parser.Entry {
	now := time.Now()
	return []parser.Entry{
		{Timestamp: now, Level: parser.LevelError, Message: "disk full"},
		{Timestamp: now, Level: parser.LevelWarn, Message: "high memory usage"},
		{Timestamp: now, Level: parser.LevelInfo, Message: "server started"},
		{Timestamp: now, Level: parser.LevelError, Message: "connection refused"},
		{Timestamp: now, Level: parser.LevelDebug, Message: "trace output"},
	}
}

func TestSplit_ByLevel(t *testing.T) {
	s := splitter.New([]splitter.Rule{
		{Name: "errors", Level: "error"},
		{Name: "warnings", Level: "warn"},
	})
	buckets := s.Split(makeEntries())

	if got := len(buckets["errors"]); got != 2 {
		t.Errorf("errors bucket: want 2, got %d", got)
	}
	if got := len(buckets["warnings"]); got != 1 {
		t.Errorf("warnings bucket: want 1, got %d", got)
	}
	if got := len(buckets["default"]); got != 2 {
		t.Errorf("default bucket: want 2, got %d", got)
	}
}

func TestSplit_ByMessageContains(t *testing.T) {
	s := splitter.New([]splitter.Rule{
		{Name: "disk", MessageContains: "disk"},
	})
	buckets := s.Split(makeEntries())

	if got := len(buckets["disk"]); got != 1 {
		t.Errorf("disk bucket: want 1, got %d", got)
	}
	if got := len(buckets["default"]); got != 4 {
		t.Errorf("default bucket: want 4, got %d", got)
	}
}

func TestSplit_FirstRuleWins(t *testing.T) {
	s := splitter.New([]splitter.Rule{
		{Name: "first", Level: "error"},
		{Name: "second", Level: "error"},
	})
	buckets := s.Split(makeEntries())

	if got := len(buckets["first"]); got != 2 {
		t.Errorf("first bucket: want 2, got %d", got)
	}
	if got := len(buckets["second"]); got != 0 {
		t.Errorf("second bucket: want 0, got %d", got)
	}
}

func TestSplit_EmptyEntries(t *testing.T) {
	s := splitter.New([]splitter.Rule{
		{Name: "errors", Level: "error"},
	})
	buckets := s.Split([]parser.Entry{})

	if len(buckets) != 0 {
		t.Errorf("expected empty buckets map, got %d keys", len(buckets))
	}
}

func TestSplit_NoRules_AllDefault(t *testing.T) {
	s := splitter.New(nil)
	entries := makeEntries()
	buckets := s.Split(entries)

	if got := len(buckets["default"]); got != len(entries) {
		t.Errorf("default bucket: want %d, got %d", len(entries), got)
	}
}
