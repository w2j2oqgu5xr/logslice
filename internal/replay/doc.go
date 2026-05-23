// Package replay provides time-accurate replay of historical log streams.
//
// A Replayer reads log entries from any replay.Source and re-emits them over
// a channel, inserting delays proportional to the original inter-entry gaps
// scaled by a configurable speed multiplier.
//
// Basic usage:
//
//	// Replay a file at double speed.
//	 src := reader.NewFileSource("archive.log")
//	 r := replay.New(src, 2.0)
//	 entries, errs := r.Run()
//	 for e := range entries {
//	     fmt.Println(e.Message)
//	 }
//	 if err := <-errs; err != nil {
//	     log.Fatal(err)
//	 }
//
// Speed multiplier semantics:
//   - 0   — no delay, emit as fast as possible
//   - 1.0 — real-time replay
//   - N   — N× faster than real time
package replay
