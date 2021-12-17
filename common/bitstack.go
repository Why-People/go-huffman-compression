package common

import (
	"sync"

	"github.com/dropbox/godropbox/container/bitvector"
)

// A thread safe wrapper around a Bit Vector in the form of a stack

type BitStack interface {
	Push(bit byte)
	Pop() bool
	Peek() bool
	Size() int
}

type threadSafeBs struct {
	mutex *sync.Mutex
	vec *bitvector.BitVector
	top int
}

func (b *threadSafeBs) Push(bit byte) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.vec.Append(bit)
	b.top++
}

func (b *threadSafeBs) Pop(i int) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.top--
	x := b.vec.Element(b.top) >= 1
	b.vec.Set(0, b.top)
	return x
}

func (b *threadSafeBs) Peek() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.vec.Element(b.top - 1) >= 1
}

func (b *threadSafeBs) Size() int {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.top
}

// Creates a new 8 Bit Bit Stack
// capacity: The number of bits to allocate.
func NewBitStack() *threadSafeBs {
	return &threadSafeBs {
		mutex: new(sync.Mutex),
		vec:   bitvector.NewBitVector([]byte{}, 1),
	}
}