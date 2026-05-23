package replay

import (
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// Source defines a source of log entries for replay.
type Source interface {
	Lines() (<-chan string, <-chan error)
}

// Replayer replays log entries at a scaled rate relative to their original
// timestamps, allowing time-accurate simulation of historical log streams.
type Replayer struct {
	source    Source
	speedMult float64
	parser    func(string) (parser.Entry, error)
}

// New creates a Replayer with the given source and speed multiplier.
// A speedMult of 1.0 replays in real time; 2.0 replays at double speed;
// 0 replays with no delay (as fast as possible).
func New(source Source, speedMult float64) *Replayer {
	return &Replayer{
		source:    source,
		speedMult: speedMult,
		parser:    parser.ParseLine,
	}
}

// Run emits parsed entries over the returned channel, introducing delays
// proportional to the original inter-entry time gaps scaled by speedMult.
func (r *Replayer) Run() (<-chan parser.Entry, <-chan error) {
	out := make(chan parser.Entry)
	errs := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errs)

		lines, lineErrs := r.source.Lines()
		var prev time.Time

		for line := range lines {
			entry, err := r.parser(line)
			if err != nil {
				continue
			}

			if r.speedMult > 0 && !prev.IsZero() {
				gap := entry.Timestamp.Sub(prev)
				if gap > 0 {
					scaled := time.Duration(float64(gap) / r.speedMult)
					time.Sleep(scaled)
				}
			}

			prev = entry.Timestamp
			out <- entry
		}

		if err, ok := <-lineErrs; ok && err != nil {
			errs <- err
		}
	}()

	return out, errs
}
