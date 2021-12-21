package main

import (
	"fmt"
	"os"

	"io.whypeople/huffman/common"
	"io.whypeople/huffman/compress"

	"github.com/akamensky/argparse"
)

const OUT_FLAGS = os.O_CREATE | os.O_WRONLY

func main() {
	// Argument parsing
	argparser := argparse.NewParser("huffman", "A simple Huffman Encoder/Decoder written for educational purposes.")


	// File args
	infileOpts := &argparse.Options{Required: true, Help: "Required Input File"}
	infile := argparser.File("i", "infile", os.O_RDWR, 0600, infileOpts)
	outfileOpts := &argparse.Options{Required: true, Help: "Required Output File Path"}
	outfile := argparser.File("o", "outfile", OUT_FLAGS, 0600, outfileOpts)
	defer outfile.Close()
	defer infile.Close()
	
	// Mode
	decodeOpts := &argparse.Options{Required: false, Help: "Decode Mode"}
	decode := argparser.Flag("d", "decode", decodeOpts)
	encodeOPts := &argparse.Options{Required: false, Help: "Encode Mode"}
	encode := argparser.Flag("e", "encode", encodeOPts)

	// Concurrency
	goroutineOpts := &argparse.Options{Required: false, Help: "Maximum Number of Goroutines to use", Default: 4}
	goroutines := argparser.Int("g", "goroutines", goroutineOpts)
	
	// Parse args
	err := argparser.Parse(os.Args)
	if err != nil {
		fmt.Println(argparser.Usage(err))
		return
	}

	// Handle mode args
	if *decode && *encode {
		fmt.Println(argparser.Usage("Must specify at most 1 mode flag (-e or -d)"))
		return
	}

	if !*decode && !*encode {
		fmt.Println(argparser.Usage("Must specify at least 1 mode flag (-e or -d)"))
		return
	}

	// Handle concurrency args
	if *goroutines < 1 {
		fmt.Println(argparser.Usage("Must specify at least 1 goroutine"))
		return
	}

	if *encode {
		fi, err := compress.CompressFile(infile, outfile, *goroutines)
		if err != nil {
			panic(err.Error())
		}
		inSize := common.GetFileSize(infile)
		outSize := common.GetFileSize(fi)
		fmt.Println("Uncompressed File Size:", inSize)
		fmt.Println("Compressed file size:", outSize)
		fmt.Printf("Compression Ratio: %3.2v\n", float64(inSize) / float64(outSize))
		fmt.Printf("Space Saving: %3.4v%%\n", (float64(inSize - outSize) / float64(inSize)) * 100.0)
	} else {
		fmt.Println("Decode Mode")
	}
}
