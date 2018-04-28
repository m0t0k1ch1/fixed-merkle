package merkle

import (
	"bytes"
	"errors"
)

var (
	ErrTooManyLeaves       = errors.New("number of leaves exceeds upper limit")
	ErrLeafIndexOutOfRange = errors.New("leaf index is out of range")
)

type Tree struct {
	config *Config
	nodes  []*Node
	levels [][]*Node
}

func NewTree(conf *Config, leaves [][]byte, hashed bool) (*Tree, error) {
	if len(leaves) > conf.allLeavesNum {
		return nil, ErrTooManyLeaves
	}

	if !hashed {
		for i, leaf := range leaves {
			conf.hasher.Reset()
			if _, err := conf.hasher.Write(leaf); err != nil {
				return nil, err
			}
			leaves[i] = conf.hasher.Sum(nil)
		}
	}

	tree := &Tree{
		config: conf,
	}
	tree.setUpBaseNodes(leaves)

	if err := tree.build(); err != nil {
		return nil, err
	}

	return tree, nil
}

func (tree *Tree) setUpBaseNodes(leaves [][]byte) {
	conf := tree.config
	leavesNum := len(leaves)

	tree.nodes = make([]*Node, conf.allNodesNum)
	tree.levels = make([][]*Node, conf.depth+1)
	tree.levels[conf.depth] = make([]*Node, conf.allLeavesNum)

	for i := 0; i < leavesNum; i++ {
		node := newNode(leaves[i], nil, nil)
		tree.nodes[i] = node
		tree.levels[conf.depth][i] = node
	}

	for i := leavesNum; i < conf.allLeavesNum; i++ {
		node := newNode(make([]byte, conf.hashSize, conf.hashSize), nil, nil)
		tree.nodes[i] = node
		tree.levels[conf.depth][i] = node
	}
}

func (tree *Tree) build() error {
	conf := tree.config

	for d := conf.depth; d > 0; d-- {
		level := tree.levels[d]

		nextDepth := d - 1
		tree.levels[nextDepth] = make([]*Node, len(level)/2)

		for i := 0; i < len(level); i += 2 {
			left := level[i]
			right := level[i+1]

			tree.config.hasher.Reset()
			b := concBytes(left.Bytes(), right.Bytes())
			if _, err := conf.hasher.Write(b); err != nil {
				return err
			}
			b = conf.hasher.Sum(nil)

			tree.levels[nextDepth][i/2] = newNode(b, left, right)
		}

		tree.nodes = append(tree.nodes, tree.levels[nextDepth]...)
	}

	return nil
}

func (tree *Tree) Root() *Node {
	return tree.levels[0][0]
}

func (tree *Tree) CreateMembershipProof(index int) ([]byte, error) {
	conf := tree.config

	if index < 0 || conf.allLeavesNum <= index {
		return nil, ErrLeafIndexOutOfRange
	}

	buf := bytes.NewBuffer(nil)

	for d := conf.depth; d > 0; d-- {
		var siblingIndex int
		if index%2 == 0 {
			siblingIndex = index + 1
		} else {
			siblingIndex = index - 1
		}

		siblingNode := tree.levels[d][siblingIndex]
		if _, err := buf.Write(siblingNode.Bytes()); err != nil {
			return nil, err
		}

		index /= 2
	}

	return buf.Bytes(), nil
}
