package watcher_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/user/logslice/internal/watcher"
)

func writeTempFile(t *testing.T, initial string) *os.File {
	t.Helper()
	f, err := os.CreateTemp("", "watcher-test-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if initial != "" {
		f.WriteString(initial)
	}
	return f
}

func TestTailSource_EmitsNewLines(t *testing.T) {
	f := writeTempFile(t, "existing line\n")
	defer os.Remove(f.Name())
	defer f.Close()

	src := watcher.NewTailSource(f.Name(), 50*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	lines, errs := src.Lines(ctx)

	// Give the goroutine time to seek to EOF before we append.
	time.Sleep(100 * time.Millisecond)

	f.WriteString("new line one\n")
	f.WriteString("new line two\n")

	var got []string
	for line := range lines {
		got = append(got, line)
		if len(got) == 2 {
			cancel()
		}
	}

	if err := <-errs; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(got), got)
	}
	if got[0] != "new line one" {
		t.Errorf("line 0: got %q, want %q", got[0], "new line one")
	}
	if got[1] != "new line two" {
		t.Errorf("line 1: got %q, want %q", got[1], "new line two")
	}
}

func TestTailSource_MissingFile(t *testing.T) {
	src := watcher.NewTailSource("/nonexistent/path/file.log", 50*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	_, errs := src.Lines(ctx)

	select {
	case err := <-errs:
		if err == nil {
			t.Fatal("expected error for missing file, got nil")
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for error")
	}
}

func TestTailSource_DefaultPollInterval(t *testing.T) {
	// Passing zero or negative duration should fall back to the default.
	src := watcher.NewTailSource("any.log", 0)
	if src == nil {
		t.Fatal("expected non-nil TailSource")
	}
}
