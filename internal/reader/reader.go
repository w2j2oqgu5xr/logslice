// Package reader provides utilities for reading log input from various sources.
package reader

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Source represents a readable log source.
type Source interface {
	Lines() ([]string, error)
}

// FileSource reads log lines from a file on disk.
type FileSource struct {
	Path string
}

// StdinSource reads log lines from standard input.
type StdinSource struct{}

// NewFileSource creates a FileSource for the given path.
func NewFileSource(path string) *FileSource {
	return &FileSource{Path: path}
}

// NewStdinSource creates a StdinSource.
func NewStdinSource() *StdinSource {
	return &StdinSource{}
}

// Lines reads all lines from the file.
func (f *FileSource) Lines() ([]string, error) {
	file, err := os.Open(f.Path)
	if err != nil {
		return nil, fmt.Errorf("reader: open file %q: %w", f.Path, err)
	}
	defer file.Close()
	return scanLines(file)
}

// Lines reads all lines from stdin.
func (s *StdinSource) Lines() ([]string, error) {
	return scanLines(os.Stdin)
}

// scanLines reads all lines from an io.Reader.
func scanLines(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reader: scan: %w", err)
	}
	return lines, nil
}
