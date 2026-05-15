// Package runner provides the top-level orchestration layer for logslice.
//
// It bridges the gap between the parsed Config (from internal/config) and the
// processing pipeline (internal/pipeline), wiring together the reader, filter,
// and exporter components in a single Run call.
//
// Typical usage:
//
//	cfg, err := config.LoadFromFlags(os.Args[1:])
//	if err != nil {
//		log.Fatal(err)
//	}
//	r := runner.New(cfg, os.Stdout)
//	if err := r.Run(); err != nil {
//		log.Fatal(err)
//	}
//
// The runner selects the appropriate LineSource (stdin, single file, or
// multi-file) based on the InputFiles slice in Config, then delegates all
// further processing to the pipeline.
package runner
