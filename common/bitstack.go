package common

import "bytes"

// Public Interface
type BitStack interface {
	Push(bit byte)
	Pop() bool
	Peek() bool
	Size() int
	Vec() BitVec
	Copy() BitStack
	Log() string
}

// The internal struct
type bs8 struct {
	vec BitVec
	top int
}

// Pushes a bit onto the stack
// bit: The bit to push onto the stack
func (b *bs8) Push(bit byte) {
	if bit == 1 {
		b.vec.SetBit(b.top)
	} else {
		b.vec.ClrBit(b.top)
	}
	b.top++
}

// Pop returns the top bit on the stack and pops it in the process
func (b *bs8) Pop() bool {
	b.top--
	return b.vec.GetBit(b.top)
}

// Peek returns the top bit from the stack without popping it
func (b *bs8) Peek() bool {
	return b.vec.GetBit(b.top - 1)
}

// Copy returns a copy of this stack
func (b *bs8) Copy() BitStack {
	vec := NewBitStack(b.vec.Capacity())
	for i := 0; i < b.top; i++ {
		if b.vec.GetBit(i) {
			vec.Push(1)
		} else {
			vec.Push(0)
		}
	}
	return vec
}

// Size returns the size of the stack
func (b *bs8) Size() int {
	return b.top
}

// Vec returns a copy of this stack's Bit Vector
func (b *bs8) Vec() BitVec {
	return b.vec
}

// Log returns a string representation of the stack suitable for logging
func (b *bs8) Log() string {
	buffer := new(bytes.Buffer)
	for i := 0; i < b.top; i++ {
		bit := b.vec.GetBit(i)
		if bit {
			buffer.WriteString("1")
		} else {
			buffer.WriteString("0")
		}
	}
	return buffer.String()
}

// NewBitStack returns a new BitStack with the given size
// size: The size of the stack in bits
func NewBitStack(size uint64) BitStack {
	return &bs8{vec: NewVector(size), top: 0}
}