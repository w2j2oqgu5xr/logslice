// Package filter provides log entry filtering for logslice.
//
// It allows callers to narrow down a slice of parsed log entries based on
// one or more criteria:
//
//   - Level       — exact (case-insensitive) match on the log level
//                   (e.g. "INFO", "WARN", "ERROR").
//   - StartTime   — only entries whose timestamp is on or after this value.
//   - EndTime     — only entries whose timestamp is on or before this value.
//   - MsgContains — case-insensitive substring search on the message field.
//
// Criteria are combined with logical AND; an empty/zero value for any field
// means that criterion is skipped.
//
// Example usage:
//
//	result := filter.Filter(entries, filter.Options{
//		Level:       "ERROR",
//		MsgContains: "timeout",
//	})
package filter
