package compress

import (
	"errors"
	"fmt"
	"io"
	"os"

	"io.whypeople/huffman/common"
)

// Return type for compressed file
type HuffCompressedFile struct {
	file *os.File
	header *common.HuffHeader
	err error
}

// CompressFile returns a data type with information about the compressed file
// infile: The file to be compressed
// outfile: The file to write the compressed data to
func CompressFile(infile *os.File, outfile *os.File, maxGoroutines int) HuffCompressedFile {
	compressedFile := HuffCompressedFile{outfile, nil, nil}

	// Make sure file pointers are valid
	if infile == nil || outfile == nil {
		compressedFile.err = errors.New("infile and outfile cannot be nil")
		return compressedFile
	}

	// Build Histogram
	// TODO: Make this concurrent

	histogram := NewHistogram()
	buf := make([]byte, common.READ_BLOCK_SIZE)

	for {
		// Read from file
		total, err := infile.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			compressedFile.err = err
			return compressedFile
		}

		// Increment the weight for each byte
		for _, b := range buf[:total] {
			histogram.IncrementWeight(b)
		}
	}
	

	for i := 0; i < common.ALPHABET_SIZE; i++ {
		w := histogram.GetWeight(byte(i))
		if w > 0 {
			fmt.Println(string(rune(i)), w)
		}
	}

	return compressedFile
}