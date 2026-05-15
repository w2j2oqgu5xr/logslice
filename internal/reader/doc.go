// Package reader provides abstractions for reading log input from one or more
// sources, including files on disk and standard input.
//
// # Sources
//
// A Source is any type that can produce a slice of raw log lines. The package
// ships with three concrete implementations:
//
//   - FileSource — reads from a single file path.
//   - StdinSource — reads from os.Stdin.
//   - MultiSource — fans in multiple Sources, concatenating their lines in
//     order. Useful when processing several log files in a single run.
//
// Empty lines are automatically discarded by all sources so that downstream
// parsers never receive blank input.
//
// # Example
//
//	src := reader.NewFileSource("/var/log/app.log")
//	lines, err := src.Lines()
//	if err != nil {
//		log.Fatal(err)
//	}
package reader
