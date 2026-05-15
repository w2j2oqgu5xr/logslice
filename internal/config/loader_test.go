package config

import (
	"flag"
	"testing"
	"time"
)

func newFS() *flag.FlagSet {
	return flag.NewFlagSet("test", flag.ContinueOnError)
}

func TestLoadFromFlags_Defaults(t *testing.T) {
	cfg, err := LoadFromFlags(newFS(), []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputFormat != FormatJSON {
		t.Errorf("expected JSON format, got %q", cfg.OutputFormat)
	}
}

func TestLoadFromFlags_LevelAndMessage(t *testing.T) {
	cfg, err := LoadFromFlags(newFS(), []string{"-level=ERROR", "-message=timeout"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Level != "ERROR" {
		t.Errorf("expected level ERROR, got %q", cfg.Level)
	}
	if cfg.MessageQ != "timeout" {
		t.Errorf("expected message 'timeout', got %q", cfg.MessageQ)
	}
}

func TestLoadFromFlags_CSVFormat(t *testing.T) {
	cfg, err := LoadFromFlags(newFS(), []string{"-format=csv"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputFormat != FormatCSV {
		t.Errorf("expected CSV format, got %q", cfg.OutputFormat)
	}
}

func TestLoadFromFlags_WindowAggregate(t *testing.T) {
	cfg, err := LoadFromFlags(newFS(), []string{"-aggregate=window", "-window=60"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AggregateBy != "window" {
		t.Errorf("expected aggregate 'window', got %q", cfg.AggregateBy)
	}
	if cfg.WindowSize != 60*time.Second {
		t.Errorf("expected 60s window, got %v", cfg.WindowSize)
	}
}

func TestLoadFromFlags_InvalidFormat(t *testing.T) {
	_, err := LoadFromFlags(newFS(), []string{"-format=toml"})
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestLoadFromFlags_InvalidStartTime(t *testing.T) {
	_, err := LoadFromFlags(newFS(), []string{"-start=not-a-time"})
	if err == nil {
		t.Error("expected error for invalid start time")
	}
}

func TestLoadFromFlags_InputFiles(t *testing.T) {
	cfg, err := LoadFromFlags(newFS(), []string{"file1.log", "file2.log"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.InputFiles) != 2 {
		t.Errorf("expected 2 input files, got %d", len(cfg.InputFiles))
	}
}
