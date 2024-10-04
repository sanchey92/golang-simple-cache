package storage

import "sync"

// Storage is a simple thread-safe in-memory key-value store.
type Storage struct {
	mu   sync.RWMutex
	Data map[string][]byte
}

func New() *Storage {
	return &Storage{
		Data: make(map[string][]byte),
	}
}

// Set stores the provided key-value pair in the storage.
//
// Params:
// - key ([]byte): The key under which the value will be stored.
// - val ([]byte): The value to be stored.
//
// Returns:
// - error: Always returns nil for this simple implementation, but can be extended for error handling.
func (s *Storage) Set(key, val []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Data[string(key)] = val
	return nil
}

// Get retrieves the value for the given key from the storage.
//
// Params:
// - key ([]byte): The key whose value should be retrieved.
//
// Returns:
// - []byte: The value associated with the key.
// - bool: A flag indicating if the key was found (true) or not (false).
func (s *Storage) Get(key []byte) ([]byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.Data[string(key)]
	return val, ok
}
