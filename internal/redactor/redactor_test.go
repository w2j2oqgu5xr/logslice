package redactor_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
	"github.com/logslice/logslice/internal/redactor"
)

func makeEntry(msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     parser.LevelInfo,
		Message:   msg,
	}
}

func TestRedact_Password(t *testing.T) {
	r := redactor.NewDefault()
	e := r.Redact(makeEntry("login failed password=hunter2 for user"))
	if contains(e.Message, "hunter2") {
		t.Errorf("expected password to be redacted, got: %s", e.Message)
	}
}

func TestRedact_Email(t *testing.T) {
	r := redactor.NewDefault()
	e := r.Redact(makeEntry("user alice@example.com logged in"))
	if contains(e.Message, "alice@example.com") {
		t.Errorf("expected email to be redacted, got: %s", e.Message)
	}
}

func TestRedact_CreditCard(t *testing.T) {
	r := redactor.NewDefault()
	e := r.Redact(makeEntry("payment with card 4111111111111111 approved"))
	if contains(e.Message, "4111111111111111") {
		t.Errorf("expected card number to be redacted, got: %s", e.Message)
	}
}

func TestRedact_NoSensitiveData(t *testing.T) {
	r := redactor.NewDefault()
	original := "server started on port 8080"
	e := r.Redact(makeEntry(original))
	if e.Message != original {
		t.Errorf("expected message unchanged, got: %s", e.Message)
	}
}

func TestApply_RedactsAll(t *testing.T) {
	r := redactor.NewDefault()
	entries := []parser.Entry{
		makeEntry("token=abc123 used"),
		makeEntry("no secrets here"),
		makeEntry("api_key=xyz987 called"),
	}
	out := r.Apply(entries)
	if contains(out[0].Message, "abc123") {
		t.Errorf("entry 0 token not redacted: %s", out[0].Message)
	}
	if out[1].Message != "no secrets here" {
		t.Errorf("entry 1 unexpectedly modified: %s", out[1].Message)
	}
	if contains(out[2].Message, "xyz987") {
		t.Errorf("entry 2 api_key not redacted: %s", out[2].Message)
	}
}

func TestNew_CustomRule(t *testing.T) {
	import_re := mustCompile(`\bSSN:\d{3}-\d{2}-\d{4}\b`)
	r := redactor.New([]redactor.Rule{
		{Pattern: import_re, Replacement: "SSN:[REDACTED]"},
	})
	e := r.Redact(makeEntry("user SSN:123-45-6789 registered"))
	if contains(e.Message, "123-45-6789") {
		t.Errorf("expected SSN redacted, got: %s", e.Message)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}

func mustCompile(pattern string) *regexp.Regexp {
	import "regexp"
	return regexp.MustCompile(pattern)
}
