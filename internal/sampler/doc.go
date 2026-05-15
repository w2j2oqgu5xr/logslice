// Package sampler provides log entry sampling strategies for use in
// high-throughput log pipelines where processing every entry is unnecessary
// or cost-prohibitive.
//
// Two strategies are supported:
//
//   - StrategyRandom – keeps each entry independently with probability Rate
//     (a float64 in the range [0.0, 1.0]).
//
//   - StrategyNth – keeps every N-th entry in the order they are received,
//     using an atomic counter so the sampler is safe for concurrent use.
//
// Example usage:
//
//	s := sampler.New(sampler.Options{
//		Strategy: sampler.StrategyNth,
//		N:        10,
//	})
//	sampled := s.Sample(entries)
//
// Invalid option values are clamped to safe defaults rather than returning
// an error, keeping call sites simple.
package sampler
