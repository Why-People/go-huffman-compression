package common

const BITS = 8

// The Public Interface
type BitVec interface {
	SetBit(i int)
	ClrBit(i int)
	GetBit(i int) bool
	RawData() []byte
	Copy() BitVec
	Capacity() uint64
}

// Internal Data Struct
type bv8 struct {
	data []byte
	cap  uint64
}

// Sets a bit in the vector to 1
// i: The index of the bit to set
func (bv *bv8) SetBit(i int) {
	bv.data[i / BITS] |= (1 << uint(i % BITS))
}

// GetBit returns the bit at the specified index
// i: The index of the bit to get
func (bv *bv8) GetBit(i int) bool {
	return bv.data[i / BITS] & (1 << uint(i % BITS)) != 0
}

// Sets a bit in the vector to 0
// i: The index of the bit to set
func (bv *bv8) ClrBit(i int) {
	bv.data[i / BITS] &= ^(1 << uint(i % BITS))
}

// RawData returns the raw byte data from this vector
func (bv *bv8) RawData() []byte {
	return bv.data
}

// Capacity returns the capacity of the vector in bits
func (bv *bv8) Capacity() uint64 {
	return bv.cap
}

// Copy returns a copy of this vector
func (bv *bv8) Copy() BitVec {
	return &bv8{data: bv.data, cap: bv.cap}
}

// NewVectorFromData creates a BitVector from data
// data: The orginal data
// cap: The capacity of the vector
func NewVectorFromData(data []byte, cap uint64) BitVec {
	return &bv8{data: data, cap: cap}
}

// NewVector creates a BitVector with the specified capacity
// cap: The capacity of the vector in bits
func NewVector(size uint64) BitVec {
	return &bv8{data: make([]byte, size / BITS), cap: size}
}