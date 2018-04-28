package merkle

import (
	"bytes"
	"errors"
)

var (
	ErrLeafIndexOutOfRange = errors.New("leaf index is out of range")
)

type Tree struct {
	Nodes  []*Node
	Levels [][]*Node
}

func (tree *Tree) Root() *Node {
	return tree.Levels[0][0]
}

func (tree *Tree) CreateMembershipProof(index int) ([]byte, error) {
	depth := len(tree.Levels) - 1
	leavesNum := len(tree.Levels[depth])

	if index < 0 || leavesNum <= index {
		return nil, ErrLeafIndexOutOfRange
	}

	buf := bytes.NewBuffer(nil)

	for i := depth; i > 0; i-- {
		var siblingIndex int
		if index%2 == 0 {
			siblingIndex = index + 1
		} else {
			siblingIndex = index - 1
		}

		siblingNode := tree.Levels[i][siblingIndex]
		if _, err := buf.Write(siblingNode.Bytes()); err != nil {
			return nil, err
		}

		index /= 2
	}

	return buf.Bytes(), nil
}
