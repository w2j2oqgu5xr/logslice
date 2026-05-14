// Command logslice is the CLI entry point for the logslice log processor.
//
// Usage:
//
//	logslice [flags]
//
// Flags:
//
//	-format  Output format: json (default) or csv
//	-level   Filter by log level (e.g. ERROR, WARN, INFO)
//	-msg     Filter by message substring
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/pipeline"
)

func main() {
	format := flag.String("format", "json", "output format: json or csv")
	level := flag.String("level", "", "filter by log level (e.g. ERROR)")
	msg := flag.String("msg", "", "filter by message substring")
	flag.Parse()

	builder := filter.NewOptionsBuilder()
	if *level != "" {
		builder = builder.WithLevel(*level)
	}
	if *msg != "" {
		builder = builder.WithMessage(*msg)
	}

	p, err := pipeline.New(pipeline.Config{
		Input:   os.Stdin,
		Output:  os.Stdout,
		Format:  *format,
		Options: builder.Build(),
	})
	if err != nil {
		log.Fatalf("logslice: %v", err)
	}

	n, err := p.Run()
	if err != nil {
		log.Fatalf("logslice: %v", err)
	}

	fmt.Fprintf(os.Stderr, "logslice: wrote %d entries\n", n)
}
