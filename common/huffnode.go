package common

const NODE_JOIN_SYMBOL = '$'

// Huffman Tree Node

// The public interface
type HuffNode interface {
	SetLeft(left HuffNode)
	SetRight(right HuffNode)
	Data() NodeData
	Left() HuffNode
	Right() HuffNode
	IsLeaf() bool
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

// Data returns the Data From the Node (symbol and weight)
func (n *node) Data() NodeData {
	return n.data
}

// Left returns the left node
func (n *node) Left() HuffNode {
	return n.left
}

// Right returns the right node
func (n *node) Right() HuffNode {
	return n.right
}

// IsLeaf returns if the node is a leaf node
func (n *node) IsLeaf() bool {
	return n.left == nil && n.right == nil
}

// Joins returns the joining 2 nodes together
// right: The node to be joined to the right
func (n *node) Join(right HuffNode) HuffNode {
	return &node{
		data: NodeData{
			Symbol: NODE_JOIN_SYMBOL,
			Weight: n.data.Weight + right.Data().Weight,
		},
		left: n,
		right: right,
	}
}

// NewNode creates a new node
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