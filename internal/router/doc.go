// Package router provides multi-route log dispatching for logslice.
//
// A Router holds one or more named Route definitions. Each Route specifies
// its own filter.Options (level, time range, message pattern) and an export
// format ("json" or "csv"). When Dispatch is called the router applies each
// route's filters independently and writes the matching entries to the
// io.Writer provided by a WriterFactory.
//
// # Basic usage
//
//	routes := []router.Route{
//	    {Name: "errors", Options: filter.NewOptionsBuilder().WithLevel("ERROR").Build(), Format: "json"},
//	    {Name: "all",    Options: filter.Options{},                                      Format: "csv"},
//	}
//	r := router.New(routes)
//	err := r.Dispatch(entries, router.MapWriterFactory(writers))
//
// # Summary
//
// Router.Summary returns per-route aggregator.LevelCounts without performing
// any I/O, which is useful for metrics or dry-run modes.
package router
