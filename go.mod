module io.whypeople/huffman

go 1.17

require (
	github.com/akamensky/argparse v1.3.1
	io.whypeople/huffman/common v0.0.0-00010101000000-000000000000
	io.whypeople/huffman/compress v0.0.0-00010101000000-000000000000
	io.whypeople/huffman/decompress v0.0.0-00010101000000-000000000000
)

replace io.whypeople/huffman/compress => ./compress

replace io.whypeople/huffman/common => ./common

replace io.whypeople/huffman/decompress => ./decompress
