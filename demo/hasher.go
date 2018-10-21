package demo

import (
	"strconv"
)

// IDGenerator generates unique IDs.
type IDGenerator struct {
	// ch is a chanel of uint64
	ch chan uint64
}

// NewIDGenerator creates an IDGenerator.
func NewIDGenerator() *IDGenerator {
	gen := IDGenerator{
		ch: make(chan uint64),
	}

	// start a goroutine...
	go func() {
		var i uint64
		// ... that runs forever
		for {
			// writes next ID into the chanel
			i++
			gen.ch <- i
		}
	}()

	return &gen
}

// Hash returns an unique ID.
func (gen *IDGenerator) Hash(string) string {
	id := <-gen.ch
	return strconv.FormatUint(id, 10)
}
