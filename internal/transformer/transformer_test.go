package transformer_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/transformer"
)

func makeEntries(msgs ...string) []parser.Entry {
	entries := make([]parser.Entry, len(msgs))
	for i, m := range msgs {
		entries[i] = parser.Entry{
			Timestamp: time.Now(),
			Level:     parser.LevelInfo,
			Message:   m,
		}
	}
	return entries
}

func TestNormalizeMessage_TrimsAndCollapses(t *testing.T) {
	entries := makeEntries("  hello   world  ", "  single  ")
	tr := transformer.New(transformer.NormalizeMessage())
	out := tr.Apply(entries)

	if out[0].Message != "hello world" {
		t.Errorf("expected 'hello world', got %q", out[0].Message)
	}
	if out[1].Message != "single" {
		t.Errorf("expected 'single', got %q", out[1].Message)
	}
}

func TestUppercaseMessage(t *testing.T) {
	entries := makeEntries("hello world")
	tr := transformer.New(transformer.UppercaseMessage())
	out := tr.Apply(entries)

	if out[0].Message != "HELLO WORLD" {
		t.Errorf("expected 'HELLO WORLD', got %q", out[0].Message)
	}
}

func TestAddField_SetsValue(t *testing.T) {
	entries := makeEntries("test message")
	tr := transformer.New(transformer.AddField("env", "production"))
	out := tr.Apply(entries)

	if out[0].Fields["env"] != "production" {
		t.Errorf("expected field env=production, got %q", out[0].Fields["env"])
	}
}

func TestAddField_OverwritesExisting(t *testing.T) {
	e := parser.Entry{
		Timestamp: time.Now(),
		Level:     parser.LevelWarn,
		Message:   "overwrite test",
		Fields:    map[string]string{"env": "staging"},
	}
	tr := transformer.New(transformer.AddField("env", "production"))
	out := tr.Apply([]parser.Entry{e})

	if out[0].Fields["env"] != "production" {
		t.Errorf("expected overwritten value 'production', got %q", out[0].Fields["env"])
	}
}

func TestChainedTransforms(t *testing.T) {
	entries := makeEntries("  hello   world  ")
	tr := transformer.New(
		transformer.NormalizeMessage(),
		transformer.UppercaseMessage(),
		transformer.AddField("processed", "true"),
	)
	out := tr.Apply(entries)

	if out[0].Message != "HELLO WORLD" {
		t.Errorf("expected 'HELLO WORLD', got %q", out[0].Message)
	}
	if out[0].Fields["processed"] != "true" {
		t.Errorf("expected processed=true, got %q", out[0].Fields["processed"])
	}
}

func TestApply_EmptyEntries(t *testing.T) {
	tr := transformer.New(transformer.NormalizeMessage())
	out := tr.Apply([]parser.Entry{})
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}
