package common

// Huffman Tree Node

// The public interface
type HuffNode interface {
	SetLeft(left HuffNode)
	SetRight(right HuffNode)
	Data() NodeData
	Left() HuffNode
	Right() HuffNode
	Join(right HuffNode) HuffNode
}

// Simple Struct to hold the data for a node
type NodeData struct {
	Symbol byte
	Weight int
}

// Internal Node Struct
type node struct {
	data NodeData
	left HuffNode
	right HuffNode
}

// Sets the left node
// left: The left node
func (n *node) SetLeft(left HuffNode) {
	n.left = left
}

// Sets the right node
// right: The right node
func (n *node) SetRight(right HuffNode) {
	n.right = right
}

// Returns the Data From the Node
func (n *node) Data() NodeData {
	return n.data
}

// Returns the left node
func (n *node) Left() HuffNode {
	return n.left
}

// Returns the right node
func (n *node) Right() HuffNode {
	return n.right
}

// Joins 2 nodes together
// right: The node to be joined to the right
func (n *node) Join(right HuffNode) HuffNode {
	return &node{
		data: NodeData{
			Symbol: '$',
			Weight: n.data.Weight + right.Data().Weight,
		},
		left: n,
		right: right,
	}
}

// Creates a New Node
// symbol: The symbol to be stored in the node
// weight: The weight of the node
func NewNode(symbol byte, weight int) HuffNode {
	return &node{
		data: NodeData{
			Symbol: symbol,
			Weight: weight,
		},
	}
}