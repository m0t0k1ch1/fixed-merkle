package merkle

import "encoding/hex"

type Node struct {
	b     []byte
	Left  *Node
	Right *Node
}

func newNode(hashBytes []byte, left, right *Node) *Node {
	return &Node{
		b:     hashBytes,
		Left:  left,
		Right: right,
	}
}

func (node *Node) Bytes() []byte {
	return node.b
}

func (node *Node) Hex() string {
	return hex.EncodeToString(node.b)
}
