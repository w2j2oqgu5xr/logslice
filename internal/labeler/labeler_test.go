package labeler_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/labeler"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(level parser.Level, msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
	}
}

func TestLabel_MatchesByMessage(t *testing.T) {
	l := labeler.NewDefault([]labeler.Rule{
		{Label: "db", MessageContains: "database"},
	})
	e := makeEntry(parser.LevelInfo, "database connection failed")
	got := l.Label(e)
	if got.Fields["label"] != "db" {
		t.Errorf("expected label 'db', got %q", got.Fields["label"])
	}
}

func TestLabel_MatchesByLevel(t *testing.T) {
	l := labeler.NewDefault([]labeler.Rule{
		{Label: "critical", Level: parser.LevelError},
	})
	e := makeEntry(parser.LevelError, "something went wrong")
	got := l.Label(e)
	if got.Fields["label"] != "critical" {
		t.Errorf("expected label 'critical', got %q", got.Fields["label"])
	}
}

func TestLabel_NoMatch_Unchanged(t *testing.T) {
	l := labeler.NewDefault([]labeler.Rule{
		{Label: "db", MessageContains: "database"},
	})
	e := makeEntry(parser.LevelInfo, "user logged in")
	got := l.Label(e)
	if v, ok := got.Fields["label"]; ok {
		t.Errorf("expected no label, got %q", v)
	}
}

func TestLabel_FirstRuleWins(t *testing.T) {
	l := labeler.NewDefault([]labeler.Rule{
		{Label: "first", MessageContains: "timeout"},
		{Label: "second", MessageContains: "timeout"},
	})
	e := makeEntry(parser.LevelWarn, "request timeout exceeded")
	got := l.Label(e)
	if got.Fields["label"] != "first" {
		t.Errorf("expected 'first', got %q", got.Fields["label"])
	}
}

func TestLabel_CustomField(t *testing.T) {
	l := labeler.New("tag", []labeler.Rule{
		{Label: "infra", Level: parser.LevelDebug},
	})
	e := makeEntry(parser.LevelDebug, "polling service")
	got := l.Label(e)
	if got.Fields["tag"] != "infra" {
		t.Errorf("expected tag 'infra', got %q", got.Fields["tag"])
	}
}

func TestApply_LabelsAll(t *testing.T) {
	l := labeler.NewDefault([]labeler.Rule{
		{Label: "error-class", Level: parser.LevelError},
	})
	entries := []parser.Entry{
		makeEntry(parser.LevelError, "disk full"),
		makeEntry(parser.LevelError, "oom killed"),
		makeEntry(parser.LevelInfo, "startup complete"),
	}
	out := l.Apply(entries)
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
	if out[0].Fields["label"] != "error-class" {
		t.Errorf("entry 0: expected 'error-class'")
	}
	if out[1].Fields["label"] != "error-class" {
		t.Errorf("entry 1: expected 'error-class'")
	}
	if _, ok := out[2].Fields["label"]; ok {
		t.Errorf("entry 2: expected no label")
	}
}
