package router

import (
	"fmt"

	"github.com/logslice/logslice/internal/aggregator"
	"github.com/logslice/logslice/internal/config"
	"github.com/logslice/logslice/internal/exporter"
	"github.com/logslice/logslice/internal/filter"
	"github.com/logslice/logslice/internal/parser"
)

// Route describes a named output destination with its own filter and export config.
type Route struct {
	Name    string
	Options filter.Options
	Format  string
}

// Router dispatches log entries to one or more named routes based on per-route
// filter criteria and export format.
type Router struct {
	routes []Route
}

// New creates a Router with the provided routes.
func New(routes []Route) *Router {
	return &Router{routes: routes}
}

// Dispatch filters entries for each route and exports them using the given writer
// factory. writerFor receives a route name and format and must return an
// io.Writer (e.g. an open file or buffer).
func (r *Router) Dispatch(entries []parser.Entry, writerFor WriterFactory) error {
	for _, route := range r.routes {
		matched := filter.Filter(entries, route.Options)

		cfg := config.Config{
			Format: route.Format,
		}
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("route %q: invalid config: %w", route.Name, err)
		}

		w, err := writerFor(route.Name, route.Format)
		if err != nil {
			return fmt.Errorf("route %q: writer error: %w", route.Name, err)
		}

		exp := exporter.New(w)
		if err := exp.Export(matched, route.Format); err != nil {
			return fmt.Errorf("route %q: export error: %w", route.Name, err)
		}
	}
	return nil
}

// Summary returns per-route entry counts without performing I/O.
func (r *Router) Summary(entries []parser.Entry) map[string]aggregator.LevelCounts {
	result := make(map[string]aggregator.LevelCounts, len(r.routes))
	for _, route := range r.routes {
		matched := filter.Filter(entries, route.Options)
		result[route.Name] = aggregator.ByLevel(matched, route.Options)
	}
	return result
}
