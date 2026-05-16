package enricher_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/enricher"
	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     parser.LevelInfo,
		Message:   msg,
		Fields:    nil,
	}
}

func TestEnrich_AddsFields(t *testing.T) {
	e := enricher.New(map[string]string{"env": "production", "app": "logslice"})
	entry := makeEntry("hello")
	enriched := e.Enrich(entry)

	if enriched.Fields["env"] != "production" {
		t.Errorf("expected env=production, got %q", enriched.Fields["env"])
	}
	if enriched.Fields["app"] != "logslice" {
		t.Errorf("expected app=logslice, got %q", enriched.Fields["app"])
	}
}

func TestEnrich_DoesNotOverwriteExisting(t *testing.T) {
	e := enricher.New(map[string]string{"env": "production"})
	entry := makeEntry("hello")
	entry.Fields = map[string]string{"env": "staging"}

	enriched := e.Enrich(entry)
	if enriched.Fields["env"] != "staging" {
		t.Errorf("expected existing field to be preserved, got %q", enriched.Fields["env"])
	}
}

func TestEnrich_NilFieldsInitialized(t *testing.T) {
	e := enricher.New(map[string]string{"region": "us-east-1"})
	entry := makeEntry("test")

	if entry.Fields != nil {
		t.Fatal("expected nil Fields before enrichment")
	}
	enriched := e.Enrich(entry)
	if enriched.Fields == nil {
		t.Fatal("expected Fields to be initialized after enrichment")
	}
	if enriched.Fields["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %q", enriched.Fields["region"])
	}
}

func TestApply_EnrichesAll(t *testing.T) {
	e := enricher.New(map[string]string{"source": "test"})
	entries := []parser.Entry{makeEntry("a"), makeEntry("b"), makeEntry("c")}

	result := e.Apply(entries)
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	for _, entry := range result {
		if entry.Fields["source"] != "test" {
			t.Errorf("expected source=test, got %q", entry.Fields["source"])
		}
	}
}

func TestApply_EmptySlice(t *testing.T) {
	e := enricher.New(map[string]string{"env": "dev"})
	result := e.Apply([]parser.Entry{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestNewWithHostname_IncludesHostname(t *testing.T) {
	e := enricher.NewWithHostname(map[string]string{"app": "logslice"})
	entry := makeEntry("msg")
	enriched := e.Enrich(entry)

	if enriched.Fields["hostname"] == "" {
		t.Error("expected hostname to be set")
	}
	if enriched.Fields["app"] != "logslice" {
		t.Errorf("expected app=logslice, got %q", enriched.Fields["app"])
	}
}
