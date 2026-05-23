package replay_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/replay"
)

func TestDefaultOptions(t *testing.T) {
	opts := replay.DefaultOptions()
	if opts.SpeedMult != 1.0 {
		t.Errorf("expected SpeedMult=1.0, got %v", opts.SpeedMult)
	}
	if opts.MaxEntries != 0 {
		t.Errorf("expected MaxEntries=0, got %d", opts.MaxEntries)
	}
}

func TestOptionsBuilder_Chaining(t *testing.T) {
	opts := replay.NewOptionsBuilder().
		WithSpeedMult(4.0).
		WithMaxEntries(100).
		Build()

	if opts.SpeedMult != 4.0 {
		t.Errorf("expected SpeedMult=4.0, got %v", opts.SpeedMult)
	}
	if opts.MaxEntries != 100 {
		t.Errorf("expected MaxEntries=100, got %d", opts.MaxEntries)
	}
}

func TestOptionsBuilder_Defaults(t *testing.T) {
	opts := replay.NewOptionsBuilder().Build()
	def := replay.DefaultOptions()

	if opts.SpeedMult != def.SpeedMult {
		t.Errorf("SpeedMult mismatch: %v vs %v", opts.SpeedMult, def.SpeedMult)
	}
	if opts.MaxEntries != def.MaxEntries {
		t.Errorf("MaxEntries mismatch: %d vs %d", opts.MaxEntries, def.MaxEntries)
	}
}

func TestOptionsBuilder_ZeroSpeed(t *testing.T) {
	opts := replay.NewOptionsBuilder().WithSpeedMult(0).Build()
	if opts.SpeedMult != 0 {
		t.Errorf("expected SpeedMult=0, got %v", opts.SpeedMult)
	}
}
