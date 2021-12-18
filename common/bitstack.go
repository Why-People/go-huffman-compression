package common

import (
	"sync"

	"github.com/dropbox/godropbox/container/bitvector"
)

// A thread safe wrapper around a Bit Vector in the form of a stack

// Public Interface
type BitStack interface {
	Push(bit byte)
	Pop() bool
	Peek() bool
	Size() int
	Vec() BitVec
}

// Internal Struct
type threadSafeBs struct {
	mutex *sync.RWMutex
	vec *bitvector.BitVector
	top int
}

// Pushes a bit onto the stack
// bit: The bit to push onto the stack
func (b *threadSafeBs) Push(bit byte) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.vec.Append(bit)
	b.top++
}

// Pop returns the top bit on the stack and pops it in the process
func (b *threadSafeBs) Pop() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.top--
	x := b.vec.Element(b.top) >= 1
	b.vec.Set(0, b.top)
	return x
}

// Peek returns the top bit from the stack
func (b *threadSafeBs) Peek() bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.vec.Element(b.top - 1) >= 1
}

// Size returns the size of the stack
func (b *threadSafeBs) Size() int {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.top
}

// Vec returns a copy of this stack's BitVector
func (b *threadSafeBs) Vec() BitVec {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return &threadSafeBv {
		mutex: new(sync.RWMutex),
		vec:   b.vec,
	}
}

// NewBitStack creates a new 8 Bit Bit Stack
// capacity: The number of bits to allocate.
func NewBitStack() *threadSafeBs {
	return &threadSafeBs {
		mutex: new(sync.RWMutex),
		vec:   bitvector.NewBitVector([]byte{}, 1),
	}
}