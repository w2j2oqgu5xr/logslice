package sampler_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/sampler"
)

func makeEntries(n int) []parser.Entry {
	entries := make([]parser.Entry, n)
	for i := range entries {
		entries[i] = parser.Entry{
			Timestamp: time.Now(),
			Level:     parser.LevelInfo,
			Message:   "test message",
		}
	}
	return entries
}

func TestSampler_RandomRate_Zero(t *testing.T) {
	s := sampler.New(sampler.Options{Strategy: sampler.StrategyRandom, Rate: 0.0})
	result := s.Sample(makeEntries(100))
	if len(result) != 0 {
		t.Errorf("expected 0 entries with rate=0, got %d", len(result))
	}
}

func TestSampler_RandomRate_One(t *testing.T) {
	s := sampler.New(sampler.Options{Strategy: sampler.StrategyRandom, Rate: 1.0})
	input := makeEntries(50)
	result := s.Sample(input)
	if len(result) != len(input) {
		t.Errorf("expected %d entries with rate=1, got %d", len(input), len(result))
	}
}

func TestSampler_NthStrategy_Basic(t *testing.T) {
	s := sampler.New(sampler.Options{Strategy: sampler.StrategyNth, N: 5})
	result := s.Sample(makeEntries(20))
	// entries 5,10,15,20 → 4 entries
	if len(result) != 4 {
		t.Errorf("expected 4 entries for N=5 over 20, got %d", len(result))
	}
}

func TestSampler_NthStrategy_N1_KeepsAll(t *testing.T) {
	s := sampler.New(sampler.Options{Strategy: sampler.StrategyNth, N: 1})
	input := makeEntries(10)
	result := s.Sample(input)
	if len(result) != len(input) {
		t.Errorf("expected all %d entries for N=1, got %d", len(input), len(result))
	}
}

func TestSampler_EmptyInput(t *testing.T) {
	s := sampler.New(sampler.Options{Strategy: sampler.StrategyRandom, Rate: 1.0})
	result := s.Sample(nil)
	if result != nil {
		t.Errorf("expected nil for empty input, got %v", result)
	}
}

func TestSampler_InvalidOptions_Clamped(t *testing.T) {
	// Rate > 1 should be clamped to 1 → all entries returned
	s := sampler.New(sampler.Options{Strategy: sampler.StrategyRandom, Rate: 99.0})
	input := makeEntries(10)
	result := s.Sample(input)
	if len(result) != len(input) {
		t.Errorf("expected all entries after rate clamp, got %d", len(result))
	}
}

func TestSampler_NthStrategy_InvalidN_Clamped(t *testing.T) {
	// N=0 should be clamped to 1 → all entries kept
	s := sampler.New(sampler.Options{Strategy: sampler.StrategyNth, N: 0})
	input := makeEntries(8)
	result := s.Sample(input)
	if len(result) != len(input) {
		t.Errorf("expected all %d entries when N clamped to 1, got %d", len(input), len(result))
	}
}
