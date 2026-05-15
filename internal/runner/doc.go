// Package runner orchestrates the full logslice pipeline by wiring together
// the reader, parser, filter, aggregator, and exporter components based on
// a resolved Config.
//
// The Runner is the top-level coordinator used by the CLI. It accepts a Config
// and an io.Writer, reads log lines from one or more sources, applies the
// configured filter options, optionally aggregates results, and writes the
// final output in the requested format (JSON or CSV).
//
// # Usage
//
//	cfg := &config.Config{
//		Files:   []string{"app.log"},
//		Level:   "error",
//		Format:  "json",
//	}
//
//	r := runner.New(cfg, os.Stdout)
//	if err := r.Run(); err != nil {
//		log.Fatal(err)
//	}
//
// When no input files are specified the runner falls back to reading from
// standard input, making it suitable for use in shell pipelines.
package runner
