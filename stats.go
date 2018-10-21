package main

import (
	"encoding/json"
	"io"
	"sync"
)

// serverStats keeps track of some server statistics.
type serverStats struct {
	mutex    sync.Mutex     // protect 'views' map from concurrent access.
	counters map[string]int // count how many times a specific URL has been viewed.
}

func newServerStats() *serverStats {
	return &serverStats{
		counters: make(map[string]int),
	}
}

// IncrRedirects increments the number of redirects for URL u.
func (s *serverStats) IncrRedirects(u string) {
	// lock the counters
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.counters[u]++
}

// Show writes the server stats as JSON into the specified writer.
func (s *serverStats) Show(w io.Writer) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// forward the writer to the JSON encoder.
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(s.counters)
}
