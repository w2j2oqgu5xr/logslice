// Package enricher provides log entry enrichment by attaching static or
// dynamic metadata fields to parsed log entries before they are exported
// or aggregated.
//
// Enrichment is non-destructive: existing fields on an entry are never
// overwritten. This makes it safe to layer multiple enrichers or to
// preserve fields that were already set during parsing.
//
// Example usage:
//
//	e := enricher.New(map[string]string{
//		"env":     "production",
//		"service": "api-gateway",
//	})
//	enrichedEntries := e.Apply(entries)
//
// To automatically include the system hostname:
//
//	e := enricher.NewWithHostname(map[string]string{"env": "staging"})
package enricher
