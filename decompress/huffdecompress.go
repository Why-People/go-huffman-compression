package decompress

import (
	"encoding/binary"
	"errors"
	"os"

	"io.whypeople/huffman/common"
)

// DecompressFile decompresses a huffman encoded file
// infile: The file to be decompressed
// outfile: The file to write the decompressed data to
// maxGoroutines: The maximum number of goroutines to use
func DecompressFile(infile *os.File, outfile *os.File, maxGoroutines int) (*os.File, error) {
	// Make sure file pointers are valid
	if infile == nil || outfile == nil {
		return nil, errors.New("infile and outfile cannot be nil")
	}

	// Read the header and validate magic number
	header := readFileHeader(infile)
	if header.MagicNumber != common.MAGIC_NUMBER {
		return nil, errors.New("invalid magic number")
	}

	// Read the tree dump
	treeDump := make([]byte, header.TreeSize)
	infile.Read(treeDump)

	// Build the huffman tree from the tree dump
	huffTreeRoot := BuildHuffmanTreeFromDump(treeDump)

	// Decompress the file
	return decompress(infile, outfile, maxGoroutines, header, huffTreeRoot)
}

// TODO: Make this concurrent
// decompress decompresses the infile and writes the decompressed data to outfile
// infile: The file to be decompressed
// outfile: The file to write the decompressed data to
// maxGoroutines: The maximum number of goroutines to use
// treeRoot: The root of the huffman tree
func decompress(infile *os.File, outfile *os.File, maxGoroutines int, header common.HuffHeader, treeRoot common.HuffNode) (*os.File, error) {
	var bitBuf common.BitVec

	readBuf := make([]byte, common.READ_BLOCK_SIZE)
	outBuf := make([]byte, common.READ_BLOCK_SIZE)
	symbolsDecoded := int64(0)
	navNode := treeRoot
	bitBufPtr := -1

	// Read the compressed data
	for symbolsDecoded < header.OriginalFileSize {
		// Read data if needed
		if bitBufPtr == -1 || bitBufPtr == common.MAX_BIT_BUFFER_SIZE {
			_, err := infile.Read(readBuf)
			if err != nil {
				return nil, err
			}
			bitBuf = common.NewVectorFromData(readBuf)
			bitBufPtr = 0
		}

		// Read the next bit and update nav node
		if bitBuf.GetBit(bitBufPtr) {
			navNode = navNode.Right()
		} else {
			navNode = navNode.Left()
		}
		bitBufPtr++

		// navNode is not a symbol if it isn't a leaf
		if !navNode.IsLeaf() {
			continue
		}

		// Write the symbol
		outBuf[symbolsDecoded % common.READ_BLOCK_SIZE] = navNode.Data().Symbol

		// Write the outbuffer if we have filled it
		if symbolsDecoded % common.READ_BLOCK_SIZE == common.READ_BLOCK_SIZE - 1 {
			_, err := outfile.Write(outBuf)
			if err != nil {
				return nil, err
			}
			// Reset outBuf
			outBuf = make([]byte, common.READ_BLOCK_SIZE)
		}

		// Reset the navigation node
		navNode = treeRoot
		symbolsDecoded++
	}

	// Write the remaining data
	_, err := outfile.Write(outBuf[:symbolsDecoded % common.READ_BLOCK_SIZE])
	if err != nil {
		return nil, err
	}

	return outfile, nil
}

// readFileHeader reads the header of a huffman encoded file
// infile: The file to read the header from
func readFileHeader(infile *os.File) common.HuffHeader {
	// Read the header
	header := common.HuffHeader{}
  binary.Read(infile, common.Endianess(), &header)
	return header
}