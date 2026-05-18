package router

import "io"

// WriterFactory is a function that returns an io.Writer for a given route name
// and format. Implementations may open files, network connections, or return
// in-memory buffers.
type WriterFactory func(name, format string) (io.Writer, error)

// MapWriterFactory builds a WriterFactory backed by a pre-populated map of
// route name → io.Writer. Useful for tests and simple CLI usage where all
// destinations are known ahead of time.
func MapWriterFactory(m map[string]io.Writer) WriterFactory {
	return func(name, _ string) (io.Writer, error) {
		w, ok := m[name]
		if !ok {
			return io.Discard, nil
		}
		return w, nil
	}
}
