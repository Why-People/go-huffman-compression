package compress

import (
	"io.whypeople/huffman/common"
)

// HistogramToHuffTree builds a Huffman Tree from a histogram and returns the root of the tree
func HistogramToHuffTree(histogram map[byte]int) common.HuffNode {
	heap := buildHuffMinHeap(histogram)

	for heap.Size() > 1 {
		left := heap.ExtractMin()
		right := heap.ExtractMin()
		node := left.Join(right)
		heap.Insert(node)
	}

	// Root node
	return heap.ExtractMin()
}

// buildHuffMinHeap builds a Huffman Min Heap from a histogram
func buildHuffMinHeap(h map[byte]int) HuffMinHeap {
	heap := NewHuffMinHeap()

	for symbol, weight := range h {
		node := common.NewNode(symbol, weight)
		heap.Insert(node)
	}

	return heap
}

// Wrapper types
type HuffCode common.BitStack
type HuffCodeTable map[byte]HuffCode 

// HuffTreeToCodeTable builds a Huffman Code Table from a Huffman Tree
// root: The root of the Huffman Tree
func HuffTreeToCodeTable(root common.HuffNode) HuffCodeTable {
	codeTable := make(HuffCodeTable)
	huffcode := common.NewBitStack(common.MAX_CODE_SIZE)
	buildHuffCodeTable(root, huffcode, codeTable)
	return codeTable
}

// Builds the Huffman Code Table recursively
func buildHuffCodeTable(n common.HuffNode, code HuffCode, codeTable HuffCodeTable) {
	if n.IsLeaf() {
		// Actual symbols are leaf nodes

		codeTable[n.Data().Symbol] = code.Copy()
	} else {
		// Traverse left
		code.Push(0)
		buildHuffCodeTable(n.Left(), code, codeTable)
		code.Pop()

		// Traverse right
		code.Push(1)
		buildHuffCodeTable(n.Right(), code, codeTable)
		code.Pop()
	}
}

// createTreeDump creates a byte array that represents the huffman tree
func CreateTreeDump(root common.HuffNode) []byte {
	dump := make([]byte, common.MAX_TREE_SIZE)
	if root == nil {
		return dump
	}

	if root.IsLeaf() {
		return []byte { 'L', root.Data().Symbol }
	}

	dump = append(dump, CreateTreeDump(root.Left())...)
	dump = append(dump, CreateTreeDump(root.Right())...)
	dump = append(dump, 'I')
	return dump
}