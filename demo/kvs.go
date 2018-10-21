package demo

import (
	"fmt"
	"sync"
)

// KVS is a very basic Key-Value Store.
type KVS struct {
	mutex  sync.Mutex
	values map[string]string
}

// NewKVS creates a new Key-Value Store.
func NewKVS() *KVS {
	return &KVS{
		values: make(map[string]string),
	}
}

// Store stores a key value pair.
func (s *KVS) Store(k, v string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.values[k] = v
	return nil
}

// Load returns the value associated with a given key.
func (s *KVS) Load(k string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	v, ok := s.values[k]
	if !ok {
		return "", fmt.Errorf("basic.KVS: unknown key '%s'", k)
	}
	return v, nil
}
