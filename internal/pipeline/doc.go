// Package pipeline provides a high-level orchestration layer for logslice.
//
// A Pipeline reads raw log lines from an [io.Reader], parses each line into a
// structured [parser.Entry], applies user-defined [filter.Options] to narrow
// the result set, and finally writes the filtered entries to an [io.Writer]
// using the [exporter] package.
//
// Basic usage:
//
//	import (
//		"os"
//		"github.com/user/logslice/internal/filter"
//		"github.com/user/logslice/internal/pipeline"
//	)
//
//	opts := filter.NewOptionsBuilder().WithLevel("ERROR").Build()
//	p, err := pipeline.New(pipeline.Config{
//		Input:   os.Stdin,
//		Output:  os.Stdout,
//		Format:  "json",
//		Options: opts,
//	})
//	if err != nil { /* handle */ }
//	n, err := p.Run()
//
// Run returns the number of entries written and any error encountered.
package pipeline
