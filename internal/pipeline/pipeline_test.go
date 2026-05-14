package pipeline_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/pipeline"
)

const sampleLogs = `2024-01-15T10:00:00Z INFO  user logged in
2024-01-15T10:01:00Z ERROR disk full
2024-01-15T10:02:00Z WARN  high memory usage
2024-01-15T10:03:00Z INFO  user logged out
invalid line without proper format
2024-01-15T10:04:00Z ERROR connection refused
`

func TestPipeline_RunJSON(t *testing.T) {
	out := &bytes.Buffer{}
	opts := filter.NewOptionsBuilder().WithLevel("ERROR").Build()

	p, err := pipeline.New(pipeline.Config{
		Input:   strings.NewReader(sampleLogs),
		Output:  out,
		Format:  "json",
		Options: opts,
	})
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	n, err := p.Run()
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 entries, got %d", n)
	}
	if !strings.Contains(out.String(), "ERROR") {
		t.Errorf("output missing ERROR entries: %s", out.String())
	}
}

func TestPipeline_RunCSV(t *testing.T) {
	out := &bytes.Buffer{}
	opts := filter.NewOptionsBuilder().Build()

	p, err := pipeline.New(pipeline.Config{
		Input:   strings.NewReader(sampleLogs),
		Output:  out,
		Format:  "csv",
		Options: opts,
	})
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	n, err := p.Run()
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if n != 5 {
		t.Errorf("expected 5 valid entries, got %d", n)
	}
}

func TestPipeline_InvalidFormat(t *testing.T) {
	_, err := pipeline.New(pipeline.Config{
		Input:  strings.NewReader(""),
		Output: &bytes.Buffer{},
		Format: "xml",
	})
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestPipeline_EmptyInput(t *testing.T) {
	out := &bytes.Buffer{}
	opts := filter.NewOptionsBuilder().Build()

	p, err := pipeline.New(pipeline.Config{
		Input:   strings.NewReader(""),
		Output:  out,
		Format:  "json",
		Options: opts,
	})
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	n, err := p.Run()
	if err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 entries for empty input, got %d", n)
	}
}
