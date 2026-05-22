package reader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logslice/internal/reader"
)

func writeTemp(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return path
}

func TestMultiSource_CombinesLines(t *testing.T) {
	dir := t.TempDir()
	p1 := writeTemp(t, dir, "a.log", "alpha\nbeta\n")
	p2 := writeTemp(t, dir, "b.log", "gamma\ndelta\n")

	m := reader.NewMultiSource(
		reader.NewFileSource(p1),
		reader.NewFileSource(p2),
	)
	lines, err := m.Lines()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
	if lines[2] != "gamma" {
		t.Errorf("expected 'gamma' at index 2, got %q", lines[2])
	}
}

func TestMultiSource_PropagatesError(t *testing.T) {
	dir := t.TempDir()
	p1 := writeTemp(t, dir, "good.log", "ok\n")

	m := reader.NewMultiSource(
		reader.NewFileSource(p1),
		reader.NewFileSource("/no/such/file.log"),
	)
	_, err := m.Lines()
	if err == nil {
		t.Error("expected error from missing source, got nil")
	}
}

func TestMultiSource_Empty(t *testing.T) {
	m := reader.NewMultiSource()
	lines, err := m.Lines()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 0 {
		t.Errorf("expected 0 lines, got %d", len(lines))
	}
}

func TestMultiSource_SingleSource(t *testing.T) {
	dir := t.TempDir()
	p1 := writeTemp(t, dir, "only.log", "one\ntwo\nthree\n")

	m := reader.NewMultiSource(reader.NewFileSource(p1))
	lines, err := m.Lines()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "one" || lines[2] != "three" {
		t.Errorf("unexpected lines: %v", lines)
	}
}
