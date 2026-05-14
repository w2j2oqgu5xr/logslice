package exporter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/exporter"
	"github.com/user/logslice/internal/parser"
)

var testTime = time.Date(2024, 5, 1, 12, 0, 0, 0, time.UTC)

var sampleEntries = []parser.Entry{
	{Timestamp: testTime, Level: "INFO", Message: "server started"},
	{Timestamp: testTime.Add(time.Minute), Level: "ERROR", Message: "connection refused"},
}

func TestExport_JSON(t *testing.T) {
	var buf bytes.Buffer
	exp := exporter.New(&buf, exporter.FormatJSON)

	if err := exp.Export(sampleEntries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "server started") {
		t.Errorf("expected message in JSON output, got: %s", out)
	}
	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected ERROR level in JSON output, got: %s", out)
	}
	if !strings.HasPrefix(strings.TrimSpace(out), "[") {
		t.Errorf("expected JSON array output, got: %s", out)
	}
}

func TestExport_CSV(t *testing.T) {
	var buf bytes.Buffer
	exp := exporter.New(&buf, exporter.FormatCSV)

	if err := exp.Export(sampleEntries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (header + 2 entries), got %d", len(lines))
	}
	if lines[0] != "timestamp,level,message" {
		t.Errorf("unexpected CSV header: %s", lines[0])
	}
	if !strings.Contains(lines[1], "INFO") {
		t.Errorf("expected INFO in first data row, got: %s", lines[1])
	}
	if !strings.Contains(lines[2], "connection refused") {
		t.Errorf("expected message in second data row, got: %s", lines[2])
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	exp := exporter.New(&buf, exporter.Format("xml"))

	if err := exp.Export(sampleEntries); err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestExport_EmptyEntries(t *testing.T) {
	var buf bytes.Buffer
	exp := exporter.New(&buf, exporter.FormatCSV)

	if err := exp.Export([]parser.Entry{}); err != nil {
		t.Fatalf("unexpected error for empty entries: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 || lines[0] != "timestamp,level,message" {
		t.Errorf("expected only header row for empty entries, got: %v", lines)
	}
}
