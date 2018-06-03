package merkle

import "encoding/hex"

type Node struct {
	b     []byte
	left  *Node
	right *Node
}

func newNode(hashBytes []byte, left, right *Node) *Node {
	return &Node{
		b:     hashBytes,
		left:  left,
		right: right,
	}
}

func (node *Node) Bytes() []byte {
	return node.b
}

func (node *Node) Hex() string {
	return hex.EncodeToString(node.b)
}

// for fmt.Stringer interface
func (node *Node) String() string {
	return "0x" + node.Hex()
}

func (node *Node) Left() *Node {
	return node.left
}

func (node *Node) Right() *Node {
	return node.right
}
