package decompress

import (
	"io.whypeople/huffman/common"
)

// BuildHuffmanTreeFromDump builds a huffman tree from a compressed file tree dump
// treeDump: The tree dump to build the tree from
func BuildHuffmanTreeFromDump(treeDump []byte) common.HuffNode {
	stack := make([]common.HuffNode, 0)

	for i := 0; i < len(treeDump); i++ {
		if treeDump[i] == common.LEAF_DUMP_CHAR {
			stack = append(stack, common.NewNode(treeDump[i + 1], 0))
			i++
		} else {
			// Pop from the stack and join nodes
			right := stack[len(stack) - 1]
			stack = stack[:len(stack) - 1]
			left := stack[len(stack) - 1]
			stack = stack[:len(stack) - 1]
			stack = append(stack, left.Join(right))
		}
	}

	// Root node
	return stack[len(stack) - 1]
}