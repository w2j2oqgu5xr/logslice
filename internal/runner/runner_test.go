package runner_test

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/config"
	"github.com/yourorg/logslice/internal/runner"
)

func writeTempLog(t *testing.T, lines string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "log*.txt")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	if _, err := f.WriteString(lines); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestRunner_RunJSON(t *testing.T) {
	path := writeTempLog(t, "2024-01-15T10:00:00Z INFO  user logged in\n2024-01-15T10:01:00Z ERROR disk full\n")

	cfg := config.Config{
		InputFiles: []string{path},
		Format:     "json",
	}
	var buf bytes.Buffer
	r := runner.New(cfg, &buf)
	if err := r.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "INFO") {
		t.Errorf("expected INFO in output, got: %s", out)
	}
	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected ERROR in output, got: %s", out)
	}
}

func TestRunner_RunWithLevelFilter(t *testing.T) {
	path := writeTempLog(t, "2024-01-15T10:00:00Z INFO  user logged in\n2024-01-15T10:01:00Z ERROR disk full\n")

	cfg := config.Config{
		InputFiles: []string{path},
		Format:     "json",
		Level:      "ERROR",
	}
	var buf bytes.Buffer
	r := runner.New(cfg, &buf)
	if err := r.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "INFO") {
		t.Errorf("expected INFO to be filtered out, got: %s", out)
	}
	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected ERROR in output, got: %s", out)
	}
}

func TestRunner_RunWithTimeRange(t *testing.T) {
	path := writeTempLog(t, "2024-01-15T09:00:00Z INFO  early\n2024-01-15T10:30:00Z WARN  in range\n2024-01-15T12:00:00Z INFO  late\n")

	since, _ := time.Parse(time.RFC3339, "2024-01-15T10:00:00Z")
	until, _ := time.Parse(time.RFC3339, "2024-01-15T11:00:00Z")

	cfg := config.Config{
		InputFiles: []string{path},
		Format:     "json",
		Since:      since,
		Until:      until,
	}
	var buf bytes.Buffer
	r := runner.New(cfg, &buf)
	if err := r.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "in range") {
		t.Errorf("expected 'in range' entry in output, got: %s", out)
	}
	if strings.Contains(out, "early") || strings.Contains(out, "late") {
		t.Errorf("expected out-of-range entries to be filtered, got: %s", out)
	}
}

func TestRunner_RunCSV(t *testing.T) {
	path := writeTempLog(t, "2024-01-15T10:00:00Z INFO  csv test\n")

	cfg := config.Config{
		InputFiles: []string{path},
		Format:     "csv",
	}
	var buf bytes.Buffer
	r := runner.New(cfg, &buf)
	if err := r.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "csv test") {
		t.Errorf("expected message in CSV output, got: %s", out)
	}
}
