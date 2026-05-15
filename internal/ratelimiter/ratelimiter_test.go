package ratelimiter_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
	"github.com/logslice/logslice/internal/ratelimiter"
)

func makeEntries(n int) []parser.Entry {
	entries := make([]parser.Entry, n)
	for i := range entries {
		entries[i] = parser.Entry{
			Level:   parser.LevelInfo,
			Message: "test message",
		}
	}
	return entries
}

func TestRateLimiter_ZeroRate_AllowsAll(t *testing.T) {
	rl := ratelimiter.New(0)
	for i := 0; i < 100; i++ {
		if !rl.Allow() {
			t.Fatalf("expected Allow() to return true with zero rate at call %d", i)
		}
	}
}

func TestRateLimiter_Apply_ZeroRate_ReturnsAll(t *testing.T) {
	rl := ratelimiter.New(0)
	entries := makeEntries(10)
	result := rl.Apply(entries)
	if len(result) != 10 {
		t.Fatalf("expected 10 entries, got %d", len(result))
	}
}

func TestRateLimiter_LimitsTokens(t *testing.T) {
	rl := ratelimiter.New(5) // 5 tokens per second
	// Immediately consume all initial tokens
	allowed := 0
	for i := 0; i < 20; i++ {
		if rl.Allow() {
			allowed++
		}
	}
	// Should allow at most 5 (initial bucket) before running dry
	if allowed > 5 {
		t.Fatalf("expected at most 5 allowed, got %d", allowed)
	}
}

func TestRateLimiter_Apply_LimitsEntries(t *testing.T) {
	rl := ratelimiter.New(3)
	entries := makeEntries(20)
	result := rl.Apply(entries)
	if len(result) > 3 {
		t.Fatalf("expected at most 3 entries, got %d", len(result))
	}
}

func TestRateLimiter_Reset_RestoresTokens(t *testing.T) {
	rl := ratelimiter.New(3)
	// Drain the bucket
	for i := 0; i < 10; i++ {
		rl.Allow()
	}
	rl.Reset()
	// After reset, should allow up to max again
	allowed := 0
	for i := 0; i < 10; i++ {
		if rl.Allow() {
			allowed++
		}
	}
	if allowed != 3 {
		t.Fatalf("expected 3 allowed after reset, got %d", allowed)
	}
}

func TestRateLimiter_RefillsOverTime(t *testing.T) {
	rl := ratelimiter.New(10)
	// Drain the bucket
	for i := 0; i < 10; i++ {
		rl.Allow()
	}
	// Simulate time passing by manipulating via sleep (short)
	time.Sleep(200 * time.Millisecond)
	// Should have refilled ~2 tokens
	allowed := 0
	for i := 0; i < 5; i++ {
		if rl.Allow() {
			allowed++
		}
	}
	if allowed < 1 {
		t.Fatal("expected at least 1 token to be refilled after 200ms")
	}
}
