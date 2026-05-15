// Package deduplicator provides log entry deduplication based on message
// content and level within a configurable time window.
package deduplicator

import (
	"sync"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// key uniquely identifies a log entry for deduplication purposes.
type key struct {
	level   parser.Level
	message string
}

// Deduplicator filters out duplicate log entries seen within a time window.
type Deduplicator struct {
	mu     sync.Mutex
	seen   map[key]time.Time
	window time.Duration
}

// New creates a new Deduplicator with the given deduplication window.
// Entries with the same level and message seen within the window are dropped.
func New(window time.Duration) *Deduplicator {
	return &Deduplicator{
		seen:   make(map[key]time.Time),
		window: window,
	}
}

// Filter returns only entries that have not been seen within the dedup window.
// It also evicts expired entries from the internal cache.
func (d *Deduplicator) Filter(entries []parser.Entry) []parser.Entry {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	d.evict(now)

	var result []parser.Entry
	for _, e := range entries {
		k := key{level: e.Level, message: e.Message}
		if last, exists := d.seen[k]; exists && now.Sub(last) < d.window {
			continue
		}
		d.seen[k] = e.Timestamp
		result = append(result, e)
	}
	return result
}

// evict removes entries from the cache that have exceeded the dedup window.
func (d *Deduplicator) evict(now time.Time) {
	for k, t := range d.seen {
		if now.Sub(t) >= d.window {
			delete(d.seen, k)
		}
	}
}

// Reset clears all tracked entries from the deduplicator cache.
func (d *Deduplicator) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[key]time.Time)
}
