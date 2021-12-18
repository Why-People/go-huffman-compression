package compress

import "sync"

// A Thread Safe Wrapper around the neccessary functions we need from a map

// The Public Interface
type Histogram interface {
	IncrementWeight(b byte)
	GetWeight(b byte) int
}

// The internal struct
type threadSafeMap struct {
	mutex *sync.RWMutex
	m map[byte]int
}

// Increments the Weight of a value in the histogram
// b: the byte to increment the weight for
func (m *threadSafeMap) IncrementWeight(b byte) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.m[b]++
}

// Gets the Weight of a value in the histogram
// b: the byte to get the weight for
func (m *threadSafeMap) GetWeight(b byte) int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.m[b]
}

// Creates a new histogram
func NewHistogram() Histogram {
	return &threadSafeMap{
		mutex: new(sync.RWMutex),
		m:     make(map[byte]int),
	}
}