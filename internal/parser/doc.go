// Package parser provides primitives for parsing structured log lines
// into Entry values that can be filtered, aggregated, and exported by
// the logslice pipeline.
//
// # Log Format
//
// The default parser expects lines in the following format:
//
//	<timestamp> [<LEVEL>] <source>: <message>
//
// Example:
//
//	2024-03-15T10:22:05 [INFO] app.server: started on port 8080
//
// Supported timestamp layouts include RFC3339, ISO 8601 without timezone,
// and common space-separated date-time strings.
//
// # Usage
//
//	entry, err := parser.ParseLine(rawLine)
//	if err != nil {
//	    // skip or log the malformed line
//	}
//	// use entry.Level, entry.Message, entry.Fields, etc.
package parser
