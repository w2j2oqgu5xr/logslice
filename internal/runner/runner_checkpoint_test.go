package runner_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/checkpoint"
)

// TestCheckpointStore_IntegrationWithRunner verifies that a checkpoint store
// correctly persists offsets across multiple simulated processing passes.
func TestCheckpointStore_IntegrationWithRunner(t *testing.T) {
	dir := t.TempDir()
	cpPath := filepath.Join(dir, "cp.json")

	store, err := checkpoint.New(cpPath)
	if err != nil {
		t.Fatalf("checkpoint.New: %v", err)
	}

	// Simulate writing a log file and tracking progress.
	logFile := filepath.Join(dir, "app.log")
	lines := []string{
		"2024-01-01T00:00:00Z INFO  first message",
		"2024-01-01T00:00:01Z WARN  second message",
		"2024-01-01T00:00:02Z ERROR third message",
	}
	content := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(logFile, []byte(content), 0o644); err != nil {
		t.Fatalf("write log: %v", err)
	}

	offset := int64(len([]byte(lines[0] + "\n")))
	if err := store.Save(logFile, offset); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Reload and verify.
	store2, err := checkpoint.New(cpPath)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	st := store2.Get(logFile)
	if st == nil {
		t.Fatal("expected state after reload")
	}
	if st.Offset != offset {
		t.Errorf("offset: want %d, got %d", offset, st.Offset)
	}

	// Simulate reading from offset.
	f, err := os.Open(logFile)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()
	if _, err := f.Seek(st.Offset, 0); err != nil {
		t.Fatalf("seek: %v", err)
	}
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(f)
	remaining := buf.String()
	for _, want := range lines[1:] {
		if !strings.Contains(remaining, want) {
			t.Errorf("expected %q in remaining content", want)
		}
	}
	if strings.Contains(remaining, lines[0]) {
		t.Error("first line should have been skipped via offset")
	}

	_ = fmt.Sprintf("offset=%d", st.Offset) // ensure fmt used
}

func TestCheckpointStore_ResetClearsState(t *testing.T) {
	dir := t.TempDir()
	store, _ := checkpoint.New(filepath.Join(dir, "cp.json"))
	_ = store.Save("some.log", 500)
	_ = store.Reset("some.log")
	if st := store.Get("some.log"); st != nil {
		t.Errorf("expected nil after reset, got %v", st)
	}
}
