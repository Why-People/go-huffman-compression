package compress

import (
	"container/heap"
	"io/whypeople/huffman/common"
)

// The min heap needed to build the Huffman Tree

// The Public Interface
type HuffMinHeap interface {
	Insert(n common.HuffNode)
	ExtractMin() common.HuffNode
	Size() int
}

// The implementation of the interface for container/heap via slices
type huffNodeHeap []common.HuffNode

// Len returns the length of the heap
func (h huffNodeHeap) Len() int {
	return len(h)
}

// Less returns true if the first node is less than the second node
// i: The index of the first node
// j: The index of the second node
func (h huffNodeHeap) Less(i, j int) bool {
	return h[i].Data().Weight < h[j].Data().Weight
}

// Swap two nodes in the heap
func (h huffNodeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Pushes a node onto the heap
func (h *huffNodeHeap) Push(x interface{}) {
	*h = append(*h, x.(common.HuffNode))
}

// Pop returns the minimum node from the heap
func (h *huffNodeHeap) Pop() interface{} {
	n := h.Len()
	x := (*h)[n-1]
	*h = (*h)[:n-1]
	return x
}

// Internal Data Struct for heap operations
type huffMinHeap struct {
	heap *huffNodeHeap
}

// Insert a node into the heap
// n: The node to be inserted
func (h *huffMinHeap) Insert(n common.HuffNode) {
	heap.Push(h.heap, n)
}

// ExtractMin returns the minimum node from the heap
func (h *huffMinHeap) ExtractMin() common.HuffNode {
	return heap.Pop(h.heap).(common.HuffNode)
}

// Size returns the size of the heap
func (h *huffMinHeap) Size() int {
	return h.heap.Len()
}

// NewHuffMinHeap returns a new min heap
func NewHuffMinHeap() HuffMinHeap {
	return &huffMinHeap{
		heap: new(huffNodeHeap),
	}
}