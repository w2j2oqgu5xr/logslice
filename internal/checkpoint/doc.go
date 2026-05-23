// Package checkpoint provides persistent offset tracking for log sources.
//
// A Store records the last processed byte offset for each named source
// (typically a file path) and persists it to a JSON file on disk. On
// restart, the runner can reload the store and resume reading from where
// it left off rather than reprocessing the entire file.
//
// Basic usage:
//
//	store, err := checkpoint.New("/var/lib/logslice/checkpoint.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Resume from last known offset.
//	if st := store.Get("app.log"); st != nil {
//		fmt.Println("resuming from offset", st.Offset)
//	}
//
//	// Persist progress after processing.
//	_ = store.Save("app.log", currentOffset)
package checkpoint
