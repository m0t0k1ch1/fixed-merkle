package merkle

import "encoding/hex"

type Node struct {
	HashBytes []byte
	Left      *Node
	Right     *Node
}

func NewNode(hashBytes []byte, left, right *Node) *Node {
	return &Node{
		HashBytes: hashBytes,
		Left:      left,
		Right:     right,
	}
}

func (node *Node) Hex() string {
	return hex.EncodeToString(node.HashBytes)
}
