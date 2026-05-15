// Package ratelimiter provides a token-bucket based rate limiter for
// controlling the throughput of log entries through the pipeline.
package ratelimiter

import (
	"sync"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

// RateLimiter controls how many log entries are allowed through per second
// using a simple token-bucket algorithm.
type RateLimiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
	clock    func() time.Time
}

// New creates a RateLimiter that allows up to ratePerSecond entries per second.
// A ratePerSecond of 0 disables rate limiting (all entries pass through).
func New(ratePerSecond float64) *RateLimiter {
	return &RateLimiter{
		tokens:   ratePerSecond,
		max:      ratePerSecond,
		rate:     ratePerSecond,
		lastTick: time.Now(),
		clock:    time.Now,
	}
}

// Allow reports whether the next log entry should be allowed through.
// It refills tokens based on elapsed time since the last call.
func (r *RateLimiter) Allow() bool {
	if r.rate == 0 {
		return true
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.clock()
	elapsed := now.Sub(r.lastTick).Seconds()
	r.lastTick = now

	r.tokens += elapsed * r.rate
	if r.tokens > r.max {
		r.tokens = r.max
	}

	if r.tokens >= 1.0 {
		r.tokens -= 1.0
		return true
	}
	return false
}

// Apply filters a slice of entries, returning only those that pass the rate limit.
func (r *RateLimiter) Apply(entries []parser.Entry) []parser.Entry {
	if r.rate == 0 {
		return entries
	}
	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		if r.Allow() {
			out = append(out, e)
		}
	}
	return out
}

// Reset restores the token bucket to its maximum capacity.
func (r *RateLimiter) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens = r.max
	r.lastTick = r.clock()
}
