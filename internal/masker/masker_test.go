package masker_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/masker"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(fields map[string]string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     parser.INFO,
		Message:   "test message",
		Fields:    fields,
	}
}

func TestMask_TargetedFieldReplaced(t *testing.T) {
	m := masker.New([]string{"token"})
	e := makeEntry(map[string]string{"token": "secret123", "user": "alice"})
	out := m.Mask(e)
	if out.Fields["token"] != "[MASKED]" {
		t.Errorf("expected [MASKED], got %q", out.Fields["token"])
	}
	if out.Fields["user"] != "alice" {
		t.Errorf("expected alice, got %q", out.Fields["user"])
	}
}

func TestMask_CaseInsensitiveFieldName(t *testing.T) {
	m := masker.New([]string{"Authorization"})
	e := makeEntry(map[string]string{"authorization": "Bearer xyz"})
	out := m.Mask(e)
	if out.Fields["authorization"] != "[MASKED]" {
		t.Errorf("expected [MASKED], got %q", out.Fields["authorization"])
	}
}

func TestMask_NoMatchingField_Unchanged(t *testing.T) {
	m := masker.New([]string{"secret"})
	e := makeEntry(map[string]string{"user": "bob"})
	out := m.Mask(e)
	if out.Fields["user"] != "bob" {
		t.Errorf("expected bob, got %q", out.Fields["user"])
	}
}

func TestMask_NilFields_Unchanged(t *testing.T) {
	m := masker.New([]string{"token"})
	e := makeEntry(nil)
	out := m.Mask(e)
	if out.Fields != nil {
		t.Errorf("expected nil fields, got %v", out.Fields)
	}
}

func TestMask_CustomPlaceholder(t *testing.T) {
	m := masker.NewWithPlaceholder([]string{"pin"}, "***")
	e := makeEntry(map[string]string{"pin": "1234"})
	out := m.Mask(e)
	if out.Fields["pin"] != "***" {
		t.Errorf("expected ***, got %q", out.Fields["pin"])
	}
}

func TestApply_MasksAll(t *testing.T) {
	m := masker.New([]string{"key"})
	entries := []parser.Entry{
		makeEntry(map[string]string{"key": "val1"}),
		makeEntry(map[string]string{"key": "val2", "other": "x"}),
	}
	out := m.Apply(entries)
	for i, e := range out {
		if e.Fields["key"] != "[MASKED]" {
			t.Errorf("entry %d: expected [MASKED], got %q", i, e.Fields["key"])
		}
	}
	if out[1].Fields["other"] != "x" {
		t.Errorf("expected other=x, got %q", out[1].Fields["other"])
	}
}

func TestApply_EmptyEntries(t *testing.T) {
	m := masker.New([]string{"key"})
	out := m.Apply([]parser.Entry{})
	if len(out) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(out))
	}
}
