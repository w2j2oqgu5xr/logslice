# logslice

A streaming log parser that filters, aggregates, and exports structured log data to JSON or CSV.

---

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git
cd logslice
go build -o logslice .
```

---

## Usage

```bash
# Filter logs by level and export to JSON
logslice --input app.log --filter level=error --output errors.json

# Aggregate logs by status code and export to CSV
logslice --input access.log --aggregate status --format csv --output report.csv

# Stream from stdin
tail -f app.log | logslice --filter level=warn --format json
```

### Flags

| Flag          | Description                              | Default  |
|---------------|------------------------------------------|----------|
| `--input`     | Path to log file (or stdin if omitted)   | stdin    |
| `--output`    | Path to output file                      | stdout   |
| `--filter`    | Key=value filter expression              | none     |
| `--aggregate` | Field to aggregate on                    | none     |
| `--format`    | Output format: `json` or `csv`           | json     |

---

## Example Output

```json
[
  { "timestamp": "2024-01-15T10:23:01Z", "level": "error", "message": "connection timeout", "service": "api" },
  { "timestamp": "2024-01-15T10:24:18Z", "level": "error", "message": "disk full", "service": "worker" }
]
```

---

## License

MIT © 2024 [yourusername](https://github.com/yourusername)