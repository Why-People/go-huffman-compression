package common

import (
	"sync"

	"github.com/dropbox/godropbox/container/bitvector"
)

// Thread Safe Wrapper around a Bit Vector

type BitVec interface {
	SetBit(i int)
	ClrBit(i int)
	GetBit(i int) bool
	RawData() []byte
}

type threadSafeBv struct {
	mutex *sync.Mutex
	vec *bitvector.BitVector
}

func (t *threadSafeBv) SetBit(i int) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.vec.Set(1, i)
}

func (t *threadSafeBv) GetBit(i int) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.vec.Element(i) >= 1
}

func (t *threadSafeBv) ClrBit(i int) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.vec.Set(0, i)
}

func (t *threadSafeBv) RawData() []byte {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.vec.Bytes()
}

// Creates a new 8 Bit Bit Vector
// capacity: The number of bits to allocate.
func NewVector(capacity int) *threadSafeBv {
	return &threadSafeBv{
		mutex: new(sync.Mutex),
		vec:   bitvector.NewBitVector([]byte{}, capacity),
	}
}