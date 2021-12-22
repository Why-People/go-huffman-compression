package common

const MAX_IO_BLOCK_SIZE = 4096
const MAGIC_NUMBER = 0xDEADBEEF
const ALPHABET_SIZE = 256
const MAX_CODE_SIZE = ALPHABET_SIZE / 8 // Bytes for a maximum, 256-bit code.
const MAX_TREE_SIZE = 3 * ALPHABET_SIZE - 1 // Maximum Huffman tree dump size.
const MAX_BIT_BUFFER_SIZE = MAX_IO_BLOCK_SIZE * 8 // Maximum size of a bit buffer.
const LEAF_DUMP_CHAR = 'L'
const INTERNAL_DUMP_CHAR = 'I'
