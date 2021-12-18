package common

// Huffman Tree Node

// The public interface
type HuffNode interface {
	SetLeft(left *node)
	SetRight(right *node)
	Data() nodeData
	Left() *node
	Right() *node
	Join(right HuffNode) *node
}

// Simple Struct to hold the data for a node
type nodeData struct {
	symbol byte
	weight int
}

// Internal Node Struct
type node struct {
	data nodeData
	left *node
	right *node
}

// Sets the left node
// left: The left node
func (n *node) SetLeft(left *node) {
	n.left = left
}

// Sets the right node
// right: The right node
func (n *node) SetRight(right *node) {
	n.right = right
}

// Returns the Data From the Node
func (n *node) Data() nodeData {
	return n.data
}

// Returns the left node
func (n *node) Left() *node {
	return n.left
}

// Returns the right node
func (n *node) Right() *node {
	return n.right
}

// Joins 2 nodes together
// right: The node to be joined to the right
func (n *node) Join(right *node) *node {
	return &node{
		data: nodeData{
			symbol: '$',
			weight: n.data.weight + right.data.weight,
		},
		left: n,
		right: right,
	}
}

// Creates a New Node
// symbol: The symbol to be stored in the node
// weight: The weight of the node
func NewNode(symbol byte, weight int) *node {
	return &node{
		data: nodeData{
			symbol: symbol,
			weight: weight,
		},
	}
}