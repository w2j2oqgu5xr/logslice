// Package redactor provides sensitive-data redaction for parsed log entries.
//
// It applies configurable regex-based rules to the message field of each
// [parser.Entry], replacing matched patterns with safe placeholder strings.
//
// # Default Rules
//
// NewDefault returns a Redactor pre-loaded with rules that scrub:
//   - Passwords and secret tokens (e.g. password=..., token=...)
//   - Email addresses
//   - Credit/debit card numbers
//
// # Custom Rules
//
// Use New with a slice of Rule values to define project-specific patterns:
//
//	r := redactor.New([]redactor.Rule{
//		{Pattern: regexp.MustCompile(`\bSSN:\d{3}-\d{2}-\d{4}\b`), Replacement: "SSN:[REDACTED]"},
//	})
//
// # Integration
//
// Redactors are typically applied after filtering and before export so that
// sensitive fields never appear in output files.
package redactor
