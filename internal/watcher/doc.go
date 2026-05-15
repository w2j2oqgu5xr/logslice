// Package watcher implements a real-time file-tail source for logslice.
//
// TailSource watches a log file for newly appended lines and streams them
// over a channel, enabling live log monitoring without reading the entire
// file from the beginning. It uses polling at a configurable interval and
// integrates with the standard context-based cancellation model used
// throughout the logslice pipeline.
//
// Basic usage:
//
//	src := watcher.NewTailSource("/var/log/app.log", 250*time.Millisecond)
//	lines, errs := src.Lines(ctx)
//	for line := range lines {
//		// process line
//	}
//	if err := <-errs; err != nil {
//		log.Fatal(err)
//	}
package watcher
