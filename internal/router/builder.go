package router

import "github.com/logslice/logslice/internal/filter"

// Builder provides a fluent API for constructing a slice of Route values
// before passing them to New.
type Builder struct {
	routes []Route
}

// NewBuilder returns an empty Builder.
func NewBuilder() *Builder {
	return &Builder{}
}

// AddRoute appends a route with the given name, filter options, and format.
func (b *Builder) AddRoute(name string, opts filter.Options, format string) *Builder {
	b.routes = append(b.routes, Route{
		Name:    name,
		Options: opts,
		Format:  format,
	})
	return b
}

// AddJSONRoute is a convenience method that adds a JSON-format route.
func (b *Builder) AddJSONRoute(name string, opts filter.Options) *Builder {
	return b.AddRoute(name, opts, "json")
}

// AddCSVRoute is a convenience method that adds a CSV-format route.
func (b *Builder) AddCSVRoute(name string, opts filter.Options) *Builder {
	return b.AddRoute(name, opts, "csv")
}

// Build returns the assembled Router.
func (b *Builder) Build() *Router {
	return New(b.routes)
}

// Routes returns the accumulated route definitions without building a Router.
func (b *Builder) Routes() []Route {
	out := make([]Route, len(b.routes))
	copy(out, b.routes)
	return out
}
