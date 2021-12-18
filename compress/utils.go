package compress

import (
	"io/whypeople/huffman/common"
)

// HistogramToHuffTree builds a Huffman Tree from a histogram and returns the root of the tree
func HistogramToHuffTree(h Histogram) common.HuffNode {
	heap := buildHuffMinHeap(h)

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
func buildHuffMinHeap(h Histogram) HuffMinHeap {
	heap := NewHuffMinHeap()

	// Insert any node that has a weight > 0
	for i := 0; i < common.ALPHABET_SIZE; i++ {
		w := h.GetWeight(byte(i))
		if w > 0 {
			node := common.NewNode(byte(i), w)
			heap.Insert(node)
		}
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
	huffcode := common.NewBitStack()
	buildHuffCodeTable(root, huffcode, codeTable)
	return codeTable
}

// Builds the Huffman Code Table recursively
func buildHuffCodeTable(n common.HuffNode, code HuffCode, codeTable HuffCodeTable) {
	if n.IsLeaf() {
		// Actual symbols are leaf nodes
		codeTable[n.Data().Symbol] = code.Copy()
	} else {
		code.Push(0)
		buildHuffCodeTable(n.Left(), code, codeTable)
		code.Pop()

		code.Push(1)
		buildHuffCodeTable(n.Right(), code, codeTable)
		code.Pop()
	}
}