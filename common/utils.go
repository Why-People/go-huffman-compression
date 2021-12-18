package common

import (
	"encoding/binary"
	"unsafe"
)

// Endianness returns the endianness of the machine
func Endianess() binary.ByteOrder {
	test := 0x01234567
	if *(*byte)(unsafe.Pointer(&test)) == 0x67 {
		return binary.LittleEndian
	} else {
		return binary.BigEndian
	}
}