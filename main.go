package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"os"
)

func main() {
	// Argument parsing
	argparser := argparse.NewParser("huffman", "A simple Huffman Encoder/Decoder written for educational purposes.")

	// File args
	inFile := argparser.File("i", "inFile", os.O_RDWR, 0600, &argparse.Options{Required: true, Help: "Required Input File"})
	outFile := argparser.File("o", "outFile", os.O_RDWR, 0600, &argparse.Options{Required: false, Help: "Output File"})
	
	// Mode
	decode := argparser.Flag("d", "decode", &argparse.Options{Required: false, Help: "Set mode to decode"})
	encode := argparser.Flag("e", "encode",  &argparse.Options{Required: false, Help: "Set mode to encode"})
	
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

	if *encode {
		fmt.Println("Encode Mode")
	} else {
		fmt.Println("Decode Mode")
	}

	fmt.Println(outFile.Name)
	fmt.Println(inFile.Name)
}
