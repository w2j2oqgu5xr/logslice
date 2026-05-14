package parser

import (
	"testing"
	"time"
)

func TestParseLine_Valid(t *testing.T) {
	line := "2024-03-15T10:22:05 [INFO] app.server: started on port 8080"
	entry, err := ParseLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Level != LevelInfo {
		t.Errorf("expected level INFO, got %q", entry.Level)
	}
	if entry.Source != "app.server" {
		t.Errorf("expected source 'app.server', got %q", entry.Source)
	}
	if entry.Message != "started on port 8080" {
		t.Errorf("unexpected message: %q", entry.Message)
	}
	if entry.Timestamp.Year() != 2024 {
		t.Errorf("unexpected timestamp year: %d", entry.Timestamp.Year())
	}
}

func TestParseLine_EmptyLine(t *testing.T) {
	_, err := ParseLine("")
	if err == nil {
		t.Error("expected error for empty line, got nil")
	}
}

func TestParseLine_InvalidFormat(t *testing.T) {
	_, err := ParseLine("this is not a valid log line")
	if err == nil {
		t.Error("expected error for invalid format, got nil")
	}
}

func TestParseLine_ErrorLevel(t *testing.T) {
	line := "2024-03-15T10:22:05 [ERROR] db.conn: connection refused"
	entry, err := ParseLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Level != LevelError {
		t.Errorf("expected level ERROR, got %q", entry.Level)
	}
}

func TestEntryIsValid(t *testing.T) {
	valid := &Entry{Timestamp: time.Now(), Message: "hello"}
	if !valid.IsValid() {
		t.Error("expected entry to be valid")
	}

	invalid := &Entry{}
	if invalid.IsValid() {
		t.Error("expected zero entry to be invalid")
	}
}

func TestEntryHasField(t *testing.T) {
	e := &Entry{Fields: map[string]string{"request_id": "abc123"}}
	if !e.HasField("request_id") {
		t.Error("expected field 'request_id' to exist")
	}
	if e.HasField("missing") {
		t.Error("expected field 'missing' to not exist")
	}
}
