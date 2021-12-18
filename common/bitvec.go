package common

import (
	"sync"

	"github.com/dropbox/godropbox/container/bitvector"
)

// Thread Safe Wrapper around a Bit Vector

// The Public Interface
type BitVec interface {
	SetBit(i int)
	ClrBit(i int)
	GetBit(i int) bool
	RawData() []byte
	Copy() BitVec
}

// Internal Data Struct
type threadSafeBv struct {
	mutex *sync.RWMutex
	vec *bitvector.BitVector
}

// Sets a bit in the vector to 1
// i: The index of the bit to set
func (t *threadSafeBv) SetBit(i int) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.vec.Set(1, i)
}

// Sets a bit in the vector to 0
// i: The index of the bit to set
func (t *threadSafeBv) ClrBit(i int) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.vec.Set(0, i)
}

// GetBit returns the bit at the specified index
// i: The index of the bit to get
func (t *threadSafeBv) GetBit(i int) bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.vec.Element(i) >= 1
}

// RawData returns the raw byte data from this vector
func (t *threadSafeBv) RawData() []byte {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.vec.Bytes()
}

// Copy returns a copy of this vector
func (t *threadSafeBv) Copy() BitVec {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return &threadSafeBv{
		mutex: new(sync.RWMutex),
		vec:   bitvector.NewBitVector(t.vec.Bytes(), t.vec.Length()),
	}
}

// NewVectorFromData creates a BitVector from data
// data: The orginal data
func NewVectorFromData(data []byte) *threadSafeBv {
	return &threadSafeBv{
		mutex: new(sync.RWMutex),
		vec:   bitvector.NewBitVector(data, len(data)),
	}
}

// NewVector creates a new empty BitVector
// capacity: The number of bits to allocate.
func NewVector(capacity int) BitVec {
	return &threadSafeBv{
		mutex: new(sync.RWMutex),
		vec:   bitvector.NewBitVector([]byte{}, capacity),
	}
}