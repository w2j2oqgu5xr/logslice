package reader

import "fmt"

// MultiSource combines multiple Sources into a single stream of lines.
type MultiSource struct {
	sources []Source
}

// NewMultiSource creates a MultiSource from the provided sources.
func NewMultiSource(sources ...Source) *MultiSource {
	return &MultiSource{sources: sources}
}

// Lines reads and concatenates lines from all underlying sources in order.
// If any source returns an error, reading stops and the error is returned.
func (m *MultiSource) Lines() ([]string, error) {
	var all []string
	for i, src := range m.sources {
		lines, err := src.Lines()
		if err != nil {
			return nil, fmt.Errorf("reader: multi source %d: %w", i, err)
		}
		all = append(all, lines...)
	}
	return all, nil
}
