module io.whypeople/huffman

go 1.17

require (
	github.com/akamensky/argparse v1.3.1
	io.whypeople/huffman/compress v0.0.0-00010101000000-000000000000
)

require (
	github.com/dropbox/godropbox v0.0.0-20200228041828-52ad444d3502 // indirect
	io.whypeople/huffman/common v0.0.0-00010101000000-000000000000 // indirect
)

replace io.whypeople/huffman/compress => ./compress

replace io.whypeople/huffman/common => ./common
