package compress

import (
	"encoding/binary"
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
	histogram, err := buildHistogramConcurrentlyFromFile(infile, maxGoroutines)
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

	for k, v := range huffCodeTable {
		fmt.Printf("%v: %v, Size: %v\n", string(k), v.Log(), v.Size())
	}

	// Dump the tree to the outfile
	treeDump := CreateTreeDump(huffTreeRoot)
	outfile.Write(treeDump)

	return compressedFile
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

// TODO: Make this more memory efficient (not having each goroutine use a local byte buffer)
// buildHistogramConcurrentlyFromFile builds a histogram from a file using the specified number of maxGoroutines
// infile: The file to build the histogram from
// maxGoroutines: The number of goroutines to use to build the histogram concurrently
func buildHistogramConcurrentlyFromFile(infile *os.File, maxGoroutines int) (map[byte]int, error) {
	histogramChan := make(chan map[byte]int)
	errChan := make(chan error)

	// Each goroutine will read from the file and build a local histogram
	for i := 0; i < maxGoroutines; i++ {
		go func() {	
			buf := make([]byte, common.READ_BLOCK_SIZE)
			localHist := make(map[byte]int)
			for {
				n, err := infile.Read(buf)
				if err != nil {
					if err == io.EOF {
						// Done reading file, send local histogram down the channel
						histogramChan <- localHist
						break
					} else {
						// Some IO Error, send it down the error channel
						errChan <- err
						return
					}
				}
				// Build local histogram
				for _, b := range buf[:n] {
					localHist[b]++
				}
			}
		}()
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