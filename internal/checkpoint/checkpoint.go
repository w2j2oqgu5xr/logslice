package checkpoint

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// State holds the persisted position within a log source.
type State struct {
	Source    string    `json:"source"`
	Offset    int64     `json:"offset"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Store persists and retrieves checkpoint state to/from a JSON file.
type Store struct {
	mu   sync.Mutex
	path string
	data map[string]*State
}

// New creates a Store backed by the given file path.
// If the file exists, existing state is loaded.
func New(path string) (*Store, error) {
	s := &Store{
		path: path,
		data: make(map[string]*State),
	}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

// Save persists the offset for the given source.
func (s *Store) Save(source string, offset int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[source] = &State{
		Source:    source,
		Offset:    offset,
		UpdatedAt: time.Now().UTC(),
	}
	return s.flush()
}

// Get returns the last saved state for a source, or nil if none exists.
func (s *Store) Get(source string) *State {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.data[source]
}

// Reset removes the checkpoint for the given source.
func (s *Store) Reset(source string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, source)
	return s.flush()
}

func (s *Store) load() error {
	f, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&s.data)
}

func (s *Store) flush() error {
	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s.data)
}
