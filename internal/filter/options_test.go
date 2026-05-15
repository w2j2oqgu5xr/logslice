package filter

import (
	"testing"
	"time"
)

func TestOptions_IsZero_Default(t *testing.T) {
	var o Options
	if !o.IsZero() {
		t.Error("expected zero Options to report IsZero() == true")
	}
}

func TestOptions_IsZero_WithLevel(t *testing.T) {
	o := Options{Level: "ERROR"}
	if o.IsZero() {
		t.Error("expected Options with Level set to report IsZero() == false")
	}
}

func TestOptions_IsZero_WithMessage(t *testing.T) {
	o := Options{MessageContains: "timeout"}
	if o.IsZero() {
		t.Error("expected Options with MessageContains set to report IsZero() == false")
	}
}

func TestOptions_IsZero_WithTimeRange(t *testing.T) {
	now := time.Now()
	o := Options{From: now}
	if o.IsZero() {
		t.Error("expected Options with From set to report IsZero() == false")
	}
}

func TestOptions_WithLevel(t *testing.T) {
	var o Options
	result := o.WithLevel("WARN")
	if result.Level != "WARN" {
		t.Errorf("expected Level=WARN, got %q", result.Level)
	}
	// original must be unchanged
	if o.Level != "" {
		t.Error("original Options should not be mutated by WithLevel")
	}
}

func TestOptions_WithTimeRange(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	var o Options
	result := o.WithTimeRange(from, to)
	if !result.From.Equal(from) {
		t.Errorf("expected From=%v, got %v", from, result.From)
	}
	if !result.To.Equal(to) {
		t.Errorf("expected To=%v, got %v", to, result.To)
	}
}

func TestOptions_WithMessage(t *testing.T) {
	var o Options
	result := o.WithMessage("database")
	if result.MessageContains != "database" {
		t.Errorf("expected MessageContains=database, got %q", result.MessageContains)
	}
	if o.MessageContains != "" {
		t.Error("original Options should not be mutated by WithMessage")
	}
}

func TestOptions_Chaining(t *testing.T) {
	from := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	o := Options{}.WithLevel("INFO").WithMessage("started").WithTimeRange(from, to)
	if o.Level != "INFO" {
		t.Errorf("expected Level=INFO, got %q", o.Level)
	}
	if o.MessageContains != "started" {
		t.Errorf("expected MessageContains=started, got %q", o.MessageContains)
	}
	if o.IsZero() {
		t.Error("chained Options should not be zero")
	}
}
