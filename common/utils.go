package common

import (
	"encoding/binary"
	"os"
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

// GetFileSize returns the size of a file in bytes
// file: the  pointer to the file
func GetFileSize(file *os.File) int64 {
	fi, err := file.Stat()
	if err != nil {
		panic(err)
	}
	return fi.Size()
}