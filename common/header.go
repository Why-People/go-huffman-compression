package common

// The header for compressed files
type HuffHeader struct {
	MagicNumber      uint32
	TreeSize         uint32
	OriginalFileSize int64
}

// CreateHeader creates a header
func CreateHeader(treeSize uint16, originalFileSize uint64) *HuffHeader {
	return &HuffHeader{MAGIC_NUMBER, uint32(treeSize), int64(originalFileSize)}
}