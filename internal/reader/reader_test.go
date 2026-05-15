package reader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logslice/internal/reader"
)

func TestFileSource_Lines_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")
	content := "line one\nline two\nline three\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	src := reader.NewFileSource(path)
	lines, err := src.Lines()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "line one" {
		t.Errorf("expected 'line one', got %q", lines[0])
	}
}

func TestFileSource_Lines_EmptyLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sparse.log")
	content := "first\n\nsecond\n\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	src := reader.NewFileSource(path)
	lines, err := src.Lines()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Errorf("expected 2 non-empty lines, got %d", len(lines))
	}
}

func TestFileSource_Lines_MissingFile(t *testing.T) {
	src := reader.NewFileSource("/nonexistent/path/to/file.log")
	_, err := src.Lines()
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestFileSource_Lines_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.log")
	if err := os.WriteFile(path, []byte(""), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	src := reader.NewFileSource(path)
	lines, err := src.Lines()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 0 {
		t.Errorf("expected 0 lines, got %d", len(lines))
	}
}
