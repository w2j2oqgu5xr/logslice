package truncator_test

import (
	"strings"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
	"github.com/logslice/logslice/internal/truncator"
)

func makeEntry(msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     parser.LevelInfo,
		Message:   msg,
	}
}

func TestTruncate_ShortMessage_Unchanged(t *testing.T) {
	tr := truncator.New(50, "...")
	e := makeEntry("short message")
	got := tr.Truncate(e)
	if got.Message != "short message" {
		t.Errorf("expected unchanged message, got %q", got.Message)
	}
}

func TestTruncate_LongMessage_Truncated(t *testing.T) {
	tr := truncator.New(10, "...")
	e := makeEntry("this is a very long message")
	got := tr.Truncate(e)
	if len(got.Message) > 10 {
		t.Errorf("expected message <= 10 bytes, got %d: %q", len(got.Message), got.Message)
	}
	if !strings.HasSuffix(got.Message, "...") {
		t.Errorf("expected suffix \"...\", got %q", got.Message)
	}
}

func TestTruncate_ExactLength_Unchanged(t *testing.T) {
	tr := truncator.New(5, "...")
	e := makeEntry("hello")
	got := tr.Truncate(e)
	if got.Message != "hello" {
		t.Errorf("expected %q, got %q", "hello", got.Message)
	}
}

func TestTruncate_ZeroMaxLen_NoTruncation(t *testing.T) {
	tr := truncator.New(0, "...")
	long := strings.Repeat("x", 500)
	e := makeEntry(long)
	got := tr.Truncate(e)
	if got.Message != long {
		t.Errorf("expected no truncation when maxLen=0")
	}
}

func TestTruncate_OriginalUnmodified(t *testing.T) {
	tr := truncator.New(10, "...")
	original := makeEntry("this is a very long message indeed")
	_ = tr.Truncate(original)
	if original.Message != "this is a very long message indeed" {
		t.Errorf("original entry should not be modified")
	}
}

func TestApply_TruncatesAll(t *testing.T) {
	tr := truncator.New(8, "...")
	entries := []parser.Entry{
		makeEntry("short"),
		makeEntry("this is definitely too long"),
		makeEntry("also very long text here"),
	}
	result := tr.Apply(entries)
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	for _, e := range result {
		if len(e.Message) > 8 {
			t.Errorf("entry message exceeds maxLen: %q", e.Message)
		}
	}
}

func TestNewDefault_Uses256Limit(t *testing.T) {
	tr := truncator.NewDefault()
	long := strings.Repeat("a", 300)
	e := makeEntry(long)
	got := tr.Truncate(e)
	if len(got.Message) > 256 {
		t.Errorf("expected <= 256 bytes, got %d", len(got.Message))
	}
	if !strings.HasSuffix(got.Message, "...") {
		t.Errorf("expected \"...\" suffix")
	}
}
