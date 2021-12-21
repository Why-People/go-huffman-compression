package compress

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"os"
	"sync/atomic"

	"io.whypeople/huffman/common"
)

// CompressFile returns a data type with information about the compressed file
// infile: The file to be compressed
// outfile: The file to write the compressed data to
func CompressFile(infile *os.File, outfile *os.File, maxGoroutines int) (*os.File, error) {

	// Make sure file pointers are valid
	if infile == nil || outfile == nil {
		return nil, errors.New("infile and outfile cannot be nil")
	}

	// Build Histogram
	histogram, err := buildHistogram(infile, maxGoroutines)
	if err != nil {
		return nil, err
	}
	
	// Build Huffman Tree
	huffTreeRoot := HistogramToHuffTree(histogram)

	// Create Header and dump it to the output file
	originalFileSize := common.GetFileSize(infile)
	writeFileHeader(outfile, len(histogram), originalFileSize)

	// If the root is null, that means the infile must've been empty
	if huffTreeRoot == nil {
		return outfile, nil
	}

	// Assign codes to each leaf (the nodes that represent bytes in infile) in the huffman tree
	huffCodeTable := HuffTreeToCodeTable(huffTreeRoot)

	// Dump the tree to the outfile
	treeDump := CreateTreeDump(huffTreeRoot)
	outfile.Write(treeDump)

	// Perform compression
	infile.Seek(0, 0)
	return compress(infile, outfile, maxGoroutines, huffCodeTable)
}

// A data struct that will be used to help keep track of the order of compressed data
type compressedBlock struct {
	data  common.BitStack // A buffer that's easy to write to bit by bit
	order int // The read order of the block
}

// compress takes a file and writes the compressed version to the output file
// infile: the file to be compressed
// outfile: the file to write the compressed data to
// codeTable: the code table to use for compression
func compress(infile *os.File, outfile *os.File, maxGoroutines int, codeTable HuffCodeTable) (*os.File, error) {
	// Create a channel to receive the compressed data from the goroutines
	compressedDataChan := make(chan compressedBlock)

	// Create a channel to receive errors from the goroutines
	errChan := make(chan error)

	// The source of truth in regards to read/write order
	orderVal := int32(-1)

	// A datastore to temporarily store the compressed data when it is finsihed being
	// processed by an eager goroutine
	totalReadsNeeded := int(math.Ceil(float64(common.GetFileSize(infile)) / float64(common.READ_BLOCK_SIZE)))
	blockStore := make([]common.BitStack, totalReadsNeeded)

	// The job each go routine will do is to read a block from the file and compress it
	compressJob := func() {
		buf := make([]byte, common.READ_BLOCK_SIZE)
		compressedBuffer := common.NewBitStack(common.READ_BLOCK_SIZE * 8) // 8 bits per byte
		for {
			// Read a block from the file
			nbytes, err := infile.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				// IO Error, send it down the error channel
				errChan <- err
				return
			}

			// Update read order val
			readOrder := atomic.AddInt32(&orderVal, 1)

			// Compress the data read
			for _, b := range buf[:nbytes] {
				// Get the code for the byte
				code := codeTable[b]
				vec := code.Vec()
				
				for i := 0; i < code.Size(); i++ {
					// Write the bit to the compressed buffer
					if vec.GetBit(i) {	
						compressedBuffer.Push(1)
					} else {
						compressedBuffer.Push(0)
					}
				}
			}

			// Write the compressed data to the channel and reset the compressed buffer
			compressedDataChan <- compressedBlock{compressedBuffer.Copy(), int(readOrder)}
			compressedBuffer.Reset()
		}
	}

	// Start the compression job for all of the workers
	for i := 0; i < maxGoroutines; i++ {
		go compressJob()
	}

	// Handle goroutine channels
	select {
	case err := <-errChan:
		return nil, err
	default:
		// Buffer used for merging compressed data blocks together
		mergedOutBuffer := common.NewBitStack(common.READ_BLOCK_SIZE * 8)
		for i := 0; i < totalReadsNeeded; i++ {
			currBlock := blockStore[i]

			// Wait for the correct block to be read if blockStore[i] is nil
			if currBlock == nil {
				// Wait for a block to be sent down the channel
				block := <-compressedDataChan
				// Incorrect block order, simply reiterate this loop until we recieve the correct block
				if i != block.order {
					i--
					// Make sure to temporarily cache this data for when it comes time to assign it to currBlock
					blockStore[block.order] = block.data
					continue
				}
				// Correct block order, assign it to currBlock
				currBlock = block.data
			}

			// Merge the current block with the merged buffer
			currBlockVec := currBlock.Vec()
			for j := 0; j < currBlock.Size(); j++ {
				// Push the bits from the current block to the merged buffer 1 by 1
				if currBlockVec.GetBit(j) {
					mergedOutBuffer.Push(1)
				} else {
					mergedOutBuffer.Push(0)
				}
				
				// Write the data in the mergedOutBuffer to the output file when it fills up
				if mergedOutBuffer.Size() == common.READ_BLOCK_SIZE * 8 {
					// Write the buffer to the output file
					outfile.Write(mergedOutBuffer.Vec().RawData())
					mergedOutBuffer.Reset()
				}
			}
			
		}

		// Write the remaining bits from the merged buffer to the outfile
		writeAmount := int(math.Ceil(float64(mergedOutBuffer.Size()) / 8))
		_, err := outfile.Write(mergedOutBuffer.Vec().RawData()[:writeAmount])
		if err != nil {
			return nil, err
		}
		return outfile, nil
	}
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