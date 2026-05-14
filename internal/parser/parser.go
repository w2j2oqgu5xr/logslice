package parser

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Common log timestamp formats to try when parsing.
var timestampFormats = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"02/Jan/2006:15:04:05 -0700",
}

// defaultPattern matches: 2006-01-02T15:04:05 [LEVEL] source: message
var defaultPattern = regexp.MustCompile(
	`^(?P<timestamp>\S+(?:\s\S+)?)\s+\[(?P<level>\w+)\]\s+(?P<source>\S+):\s+(?P<message>.+)$`,
)

// ParseLine attempts to parse a single raw log line into an Entry.
// Returns an error if the line does not match the expected format.
func ParseLine(line string) (*Entry, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("empty line")
	}

	match := defaultPattern.FindStringSubmatch(line)
	if match == nil {
		return nil, fmt.Errorf("line does not match log pattern: %q", line)
	}

	groups := make(map[string]string)
	for i, name := range defaultPattern.SubexpNames() {
		if name != "" {
			groups[name] = match[i]
		}
	}

	ts, err := parseTimestamp(groups["timestamp"])
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp %q: %w", groups["timestamp"], err)
	}

	return &Entry{
		Timestamp: ts,
		Level:     LogLevel(strings.ToUpper(groups["level"])),
		Source:    groups["source"],
		Message:   groups["message"],
		Fields:    make(map[string]string),
		Raw:       line,
	}, nil
}

func parseTimestamp(raw string) (time.Time, error) {
	for _, format := range timestampFormats {
		if t, err := time.Parse(format, raw); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("no matching timestamp format")
}
