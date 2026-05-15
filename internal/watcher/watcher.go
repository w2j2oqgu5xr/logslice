// Package watcher provides a file-watching source that emits new lines
// appended to a log file in real time, suitable for streaming pipelines.
package watcher

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// DefaultPollInterval is how often the watcher checks for new content.
const DefaultPollInterval = 250 * time.Millisecond

// TailSource watches a file and streams newly appended lines over a channel.
// It seeks to the end of the file on start and emits only new lines.
type TailSource struct {
	path         string
	pollInterval time.Duration
}

// NewTailSource creates a TailSource for the given file path.
func NewTailSource(path string, pollInterval time.Duration) *TailSource {
	if pollInterval <= 0 {
		pollInterval = DefaultPollInterval
	}
	return &TailSource{path: path, pollInterval: pollInterval}
}

// Lines opens the file, seeks to EOF, and streams new lines until ctx is done.
// Errors are sent on the returned error channel (buffered, size 1).
func (t *TailSource) Lines(ctx context.Context) (<-chan string, <-chan error) {
	lines := make(chan string)
	errs := make(chan error, 1)

	go func() {
		defer close(lines)
		defer close(errs)

		f, err := os.Open(t.path)
		if err != nil {
			errs <- err
			return
		}
		defer f.Close()

		// Seek to end so we only tail new content.
		if _, err := f.Seek(0, io.SeekEnd); err != nil {
			errs <- err
			return
		}

		reader := bufio.NewReader(f)
		ticker := time.NewTicker(t.pollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for {
					line, err := reader.ReadString('\n')
					if len(line) > 0 {
						// Strip trailing newline before sending.
						if len(line) > 0 && line[len(line)-1] == '\n' {
							line = line[:len(line)-1]
						}
						select {
						case lines <- line:
						case <-ctx.Done():
							return
						}
					}
					if err == io.EOF {
						break
					}
					if err != nil {
						errs <- err
						return
					}
				}
			}
		}
	}()

	return lines, errs
}
