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
	order int 						// The read order of the block
}

// compress takes a file and writes the compressed version to the output file
// infile: the file to be compressed
// outfile: the file to write the compressed data to
// codeTable: the code table to use for compression
func compress(infile *os.File, outfile *os.File, maxGoroutines int, codeTable HuffCodeTable) (*os.File, error) {
	totalReadsNeeded := int(math.Ceil(float64(common.GetFileSize(infile)) / float64(common.MAX_IO_BLOCK_SIZE)))

	// Create concurrent channels to handle errors and manage the order of compressed data
	compressedChannels := make([]chan common.BitStack, 0, totalReadsNeeded)
	errorChannel :=  make(chan error)
	defer close(errorChannel)
	defer func() {
		for _, c := range compressedChannels {
			close(c)
		}
	}()

	// This value will keep track of the read/write order of the compressed data
	order := int32(-1)

	// Each goroutine will read from the file and compress it
	compressJob := func() {
		buf := make([]byte, common.MAX_IO_BLOCK_SIZE)
		compressedBuffer := common.NewBitStack(common.MAX_BIT_BUFFER_SIZE)

		// Read and compress data while not at the end of the file
		for {
			nbytes, err := infile.Read(buf);
			if err == io.EOF {
				break
			}
			// Check for any IO Errors
			if err != nil {
				errorChannel <- err
				return
			}

			// Update read order val
			readOrder := atomic.AddInt32(&order, 1)

			// Compress the data read
			for _, b := range buf[:nbytes] {
				// Get the code for the byte
				code := codeTable[b]
				compressedBuffer.Append(code, 0)
			}
		
			// Copy the compressed buffer to the channel, then reset it
			compressedChannels[readOrder] <- compressedBuffer.Copy()
			compressedBuffer.Reset()
		}
	}

	// Create empty channels that hold the compressed data
	for i := 0; i < totalReadsNeeded; i++ {
		compressedChannels = append(compressedChannels, make(chan common.BitStack))
	}

	// Start the workers
	for i := 0; i < maxGoroutines; i++ {
		go compressJob()
	}

	// Merge compressed data from the workers into a single buffer
	mergedOutBuffer := common.NewBitStack(common.MAX_BIT_BUFFER_SIZE)
	select {
	case err := <-errorChannel:
		return nil, err
	default:
		for _, channel := range compressedChannels {
			// Read the data from the channel
			data := <-channel

			// Merge data from the compressed channels
			_, idx := mergedOutBuffer.Append(data, 0)
			for idx < data.Size() {
				if mergedOutBuffer.Size() >= common.MAX_BIT_BUFFER_SIZE {
					// Write the data to the output file
					outfile.Write(mergedOutBuffer.Vec().RawData())
					mergedOutBuffer.Reset()
				}
				_, idx = mergedOutBuffer.Append(data, idx)
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
	defer close(histogramChan)
	defer close(errChan)

	// The job each go routine will do is to read a block from the file and increment their own histogram
	countBytesJob := func() {
		buf := make([]byte, common.MAX_IO_BLOCK_SIZE)
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
			hist := <- histogramChan
			for k, v := range hist {
				histogram[k] += v
			}
		}
		return histogram, nil
	}
}