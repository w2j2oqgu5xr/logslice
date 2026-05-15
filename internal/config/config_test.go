package config

import (
	"testing"
	"time"
)

func TestValidate_DefaultFormat(t *testing.T) {
	c := &Config{}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.OutputFormat != FormatJSON {
		t.Errorf("expected default format JSON, got %q", c.OutputFormat)
	}
}

func TestValidate_ValidFormats(t *testing.T) {
	for _, f := range []Format{FormatJSON, FormatCSV} {
		c := &Config{OutputFormat: f}
		if err := c.Validate(); err != nil {
			t.Errorf("format %q: unexpected error: %v", f, err)
		}
	}
}

func TestValidate_UnsupportedFormat(t *testing.T) {
	c := &Config{OutputFormat: "xml"}
	if err := c.Validate(); err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestValidate_AggregateByLevel(t *testing.T) {
	c := &Config{AggregateBy: "level"}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_AggregateByWindowMissingDuration(t *testing.T) {
	c := &Config{AggregateBy: "window"}
	if err := c.Validate(); err == nil {
		t.Error("expected error when window size is zero")
	}
}

func TestValidate_AggregateByWindowValid(t *testing.T) {
	c := &Config{AggregateBy: "window", WindowSize: time.Minute}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_InvalidAggregateBy(t *testing.T) {
	c := &Config{AggregateBy: "hour"}
	if err := c.Validate(); err == nil {
		t.Error("expected error for unsupported aggregate-by")
	}
}

func TestValidate_TimeRangeInverted(t *testing.T) {
	now := time.Now()
	c := &Config{
		StartTime: now.Add(time.Hour),
		EndTime:   now,
	}
	if err := c.Validate(); err == nil {
		t.Error("expected error when end-time is before start-time")
	}
}

func TestValidate_TimeRangeValid(t *testing.T) {
	now := time.Now()
	c := &Config{
		StartTime: now,
		EndTime:   now.Add(time.Hour),
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
