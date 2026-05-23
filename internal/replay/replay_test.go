package replay_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/replay"
)

// fakeSource implements replay.Source from a slice of raw log lines.
type fakeSource struct {
	lines []string
}

func (f *fakeSource) Lines() (<-chan string, <-chan error) {
	ch := make(chan string, len(f.lines))
	errs := make(chan error, 1)
	for _, l := range f.lines {
		ch <- l
	}
	close(ch)
	close(errs)
	return ch, errs
}

func TestReplay_EmitsAllEntries(t *testing.T) {
	src := &fakeSource{
		lines: []string{
			"2024-01-01T10:00:00Z INFO first message",
			"2024-01-01T10:00:01Z INFO second message",
			"2024-01-01T10:00:02Z WARN third message",
		},
	}
	r := replay.New(src, 0) // no delay
	out, errs := r.Run()

	var entries []string
	for e := range out {
		entries = append(entries, e.Message)
	}

	if err := <-errs; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestReplay_SkipsInvalidLines(t *testing.T) {
	src := &fakeSource{
		lines: []string{
			"not a valid log line",
			"2024-01-01T10:00:00Z INFO valid",
		},
	}
	r := replay.New(src, 0)
	out, _ := r.Run()

	var count int
	for range out {
		count++
	}
	if count != 1 {
		t.Fatalf("expected 1 valid entry, got %d", count)
	}
}

func TestReplay_SpeedMultZero_NoDelay(t *testing.T) {
	src := &fakeSource{
		lines: []string{
			"2024-01-01T10:00:00Z INFO a",
			"2024-01-01T10:01:00Z INFO b", // 1 minute gap
		},
	}
	r := replay.New(src, 0)
	start := time.Now()
	out, _ := r.Run()
	for range out {
	}
	if time.Since(start) > 200*time.Millisecond {
		t.Fatal("expected near-instant replay with speedMult=0")
	}
}

func TestReplay_EmptySource(t *testing.T) {
	src := &fakeSource{lines: []string{}}
	r := replay.New(src, 1.0)
	out, errs := r.Run()

	var count int
	for range out {
		count++
	}
	if err := <-errs; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected 0 entries, got %d", count)
	}
}
