package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/checkpoint"
)

func tempPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "checkpoint.json")
}

func TestCheckpoint_SaveAndGet(t *testing.T) {
	p := tempPath(t)
	s, err := checkpoint.New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := s.Save("file.log", 42); err != nil {
		t.Fatalf("Save: %v", err)
	}
	st := s.Get("file.log")
	if st == nil {
		t.Fatal("expected state, got nil")
	}
	if st.Offset != 42 {
		t.Errorf("offset: want 42, got %d", st.Offset)
	}
	if st.Source != "file.log" {
		t.Errorf("source: want file.log, got %s", st.Source)
	}
	if st.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}
}

func TestCheckpoint_Persistence(t *testing.T) {
	p := tempPath(t)
	s1, _ := checkpoint.New(p)
	_ = s1.Save("app.log", 100)

	s2, err := checkpoint.New(p)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	st := s2.Get("app.log")
	if st == nil || st.Offset != 100 {
		t.Errorf("expected offset 100 after reload, got %v", st)
	}
}

func TestCheckpoint_GetMissing(t *testing.T) {
	s, _ := checkpoint.New(tempPath(t))
	if st := s.Get("nonexistent"); st != nil {
		t.Errorf("expected nil for unknown source, got %v", st)
	}
}

func TestCheckpoint_Reset(t *testing.T) {
	p := tempPath(t)
	s, _ := checkpoint.New(p)
	_ = s.Save("x.log", 77)
	if err := s.Reset("x.log"); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	if st := s.Get("x.log"); st != nil {
		t.Errorf("expected nil after reset, got %v", st)
	}
}

func TestCheckpoint_UpdatedAtIsRecent(t *testing.T) {
	before := time.Now().UTC().Add(-time.Second)
	s, _ := checkpoint.New(tempPath(t))
	_ = s.Save("ts.log", 1)
	st := s.Get("ts.log")
	if st.UpdatedAt.Before(before) {
		t.Errorf("UpdatedAt %v is before test start %v", st.UpdatedAt, before)
	}
}

func TestCheckpoint_MissingFileIsOK(t *testing.T) {
	p := filepath.Join(t.TempDir(), "does_not_exist.json")
	_, err := checkpoint.New(p)
	if err != nil {
		t.Errorf("expected no error for missing file, got %v", err)
	}
}

func TestCheckpoint_OverwriteOffset(t *testing.T) {
	s, _ := checkpoint.New(tempPath(t))
	_ = s.Save("log", 10)
	_ = s.Save("log", 999)
	if st := s.Get("log"); st.Offset != 999 {
		t.Errorf("expected 999, got %d", st.Offset)
	}
}

func TestCheckpoint_MultipleSourcesIsolated(t *testing.T) {
	s, _ := checkpoint.New(tempPath(t))
	_ = s.Save("a.log", 1)
	_ = s.Save("b.log", 2)
	if s.Get("a.log").Offset != 1 {
		t.Error("a.log offset mismatch")
	}
	if s.Get("b.log").Offset != 2 {
		t.Error("b.log offset mismatch")
	}
	_ = os.Remove("") // no-op, just ensure os is used
}
