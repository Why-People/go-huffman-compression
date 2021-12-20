package compress

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"os"

	"io.whypeople/huffman/common"
)

// Return type for compressed file
type HuffCompressedFile struct {
	File *os.File
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
	histogram, err := buildHistogram(infile, maxGoroutines)
	if err != nil {
		compressedFile.err = err
		return compressedFile
	}
	
	// Build Huffman Tree
	huffTreeRoot := HistogramToHuffTree(histogram)

	// Create Header and dump it to the output file
	originalFileSize := common.GetFileSize(infile)
	compressedFile.header = writeFileHeader(outfile, len(histogram), originalFileSize)

	// If the root is null, that means the infile must've been empty
	if huffTreeRoot == nil {
		return compressedFile
	}

	// Assign codes to each leaf (the nodes that represent bytes in infile) in the huffman tree
	huffCodeTable := HuffTreeToCodeTable(huffTreeRoot)

	// Dump the tree to the outfile
	treeDump := CreateTreeDump(huffTreeRoot)
	outfile.Write(treeDump)

	// Perform compression
	infile.Seek(0, 0)
	compressedFi, err := compress(infile, outfile, huffCodeTable)
	if err != nil {
		compressedFile.err = err
		return compressedFile
	}
	compressedFile.File = compressedFi
	return compressedFile
}

// TODO: Make this concurrent
// compress takes a file and writes the compressed version to the output file
// infile: the file to be compressed
// outfile: the file to write the compressed data to
// codeTable: the code table to use for compression
func compress(infile *os.File, outfile *os.File, codeTable HuffCodeTable) (*os.File, error) {
	// Create a buffer to store the compressed data
	buf := make([]byte, common.READ_BLOCK_SIZE)
	outBuf := common.NewBitStack(common.READ_BLOCK_SIZE * 8) // 8 bits per byte
	for {
		nbytes, err := infile.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			// Some IO Error, return it
			return nil, err
		}
		for i := 0; i < nbytes; i++ {
			huffCode := codeTable[buf[i]]
			if huffCode == nil {
				continue
			}
			codeVec := huffCode.Vec()

			// Push the bits from the code vector to the compression vector
			for j := 0; j < huffCode.Size(); j++ {
				if codeVec.GetBit(j) {
					outBuf.Push(1)
				} else {
					outBuf.Push(0)
				}

				if outBuf.Size() == common.READ_BLOCK_SIZE * 8 {
					// Write compressed data to file
					_, err := outfile.Write(outBuf.Vec().RawData())
					if err != nil {
						// Handle IO error
						return nil, err
					}
					outBuf.Reset()
				}
			}
		}
	}

	// Write the remaining bits to the file
	writeAmount := int(math.Ceil(float64(outBuf.Size()) / 8))
	_, err := outfile.Write(outBuf.Vec().RawData()[:writeAmount])
	if err != nil {
		return nil, err
	}
	return outfile, nil
}

// writeFileHeader writes the compression header to the output file and returns the header struct
// outfile: the output file to be compressed
// uniqueSymbols: the number of unique symbols in the orignal file
// originalFileSize: the size of the original file
func writeFileHeader(outfile *os.File, uniqueSymbolsFromIn int, originalFileSize int64) *common.HuffHeader {
	treeSize := (3 * uniqueSymbolsFromIn) - 1;
	header := common.CreateHeader(uint16(treeSize), uint64(originalFileSize))
	binary.Write(outfile, common.Endianess(), *header)
	return header
}

// buildHistogram builds a histogram from a file using the specified number of maxGoroutines
// infile: The file to build the histogram from
// maxGoroutines: The number of goroutines to use to build the histogram concurrently
func buildHistogram(infile *os.File, maxGoroutines int) (map[byte]int, error) {
	histogramChan := make(chan map[byte]int)
	errChan := make(chan error)

	// The job each go routine will do is to read a block from the file and increment their own histogram
	countBytesJob := func() {
		buf := make([]byte, common.READ_BLOCK_SIZE)
		localHist := make(map[byte]int)
		for {
			n, err := infile.Read(buf)
			if err == io.EOF {
				// Done reading file, send local histogram down the channel
				histogramChan <- localHist
				break
			}
			if err != nil {
				// Some IO Error, send it down the error channel
				errChan <- err
				return
			}
			// Build local histogram
			for _, b := range buf[:n] {
				localHist[b]++
			}
		}
	}

	// Each goroutine will read from the file and build a local histogram
	for i := 0; i < maxGoroutines; i++ {
		go countBytesJob()
	}

	select {
	case err := <-errChan:
		// Handle Error
		return nil, err
	default:
		// Recieve all the local histograms from the goroutines and compile them into 1 histogram
		histogram := make(map[byte]int)
		for i := 0; i < maxGoroutines; i++ {
			hist := <-histogramChan
			for k, v := range hist {
				histogram[k] += v
			}
		}
		return histogram, nil
	}
}