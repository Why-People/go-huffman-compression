package compress

import (
	"io/whypeople/huffman/common"
	"os"
)

// Return type for compressed file
type HuffCompressedFile struct {
	file *os.File
	header common.HuffHeader
	err error
}

// CompressFile returns a data type with information about the compressed file
// infile: The file to be compressed
// outfile: The file to write the compressed data to
func CompressFile(infile *os.File, outfile *os.File) HuffCompressedFile {
	return HuffCompressedFile{outfile, common.CreateHeader(0, 0), nil}
}