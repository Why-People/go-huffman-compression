package common

// The header for compressed files
type HuffHeader struct {
	magicNumber      uint32
	treeSize         uint16
	originalFileSize uint64
}

// CreateHeader creates a header
func CreateHeader(treeSize uint16, originalFileSize uint64) HuffHeader {
	return HuffHeader{MAGIC_NUMBER, treeSize, originalFileSize}
}