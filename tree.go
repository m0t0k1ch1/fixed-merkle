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

func NewTree(conf *Config, leaves [][]byte) (*Tree, error) {
	if uint64(len(leaves)) > conf.allLeavesNum {
		return nil, ErrTooManyLeaves
	}

	tree := &Tree{
		config: conf,
	}
	if err := tree.buildBase(leaves); err != nil {
		return nil, err
	}
	if err := tree.build(); err != nil {
		return nil, err
	}

	return tree, nil
}

func (tree *Tree) buildBase(leaves [][]byte) error {
	conf := tree.config
	leavesNum := len(leaves)

	tree.nodes = make([]*Node, conf.allNodesNum)
	tree.levels = make([][]*Node, conf.depth+1)
	tree.levels[conf.depth] = make([]*Node, conf.allLeavesNum)

	for i := 0; i < leavesNum; i++ {
		conf.hasher.Reset()
		if _, err := conf.hasher.Write(leaves[i]); err != nil {
			return err
		}
		leafHash := conf.hasher.Sum(nil)

		node := newNode(leafHash, nil, nil)
		tree.nodes[i] = node
		tree.levels[conf.depth][i] = node
	}

	emptyLeaf := make([]byte, conf.hashSize, conf.hashSize)

	conf.hasher.Reset()
	if _, err := conf.hasher.Write(emptyLeaf); err != nil {
		return err
	}
	emptyLeafHash := conf.hasher.Sum(nil)

	for i := uint64(leavesNum); i < conf.allLeavesNum; i++ {
		node := newNode(emptyLeafHash, nil, nil)
		tree.nodes[i] = node
		tree.levels[conf.depth][i] = node
	}

	return nil
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
			b := concat(left.Bytes(), right.Bytes())
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
	return tree.Level(0)[0]
}

func (tree *Tree) Level(depth uint64) []*Node {
	if tree.config.depth < depth {
		return nil
	}

	return tree.levels[depth]
}

func (tree *Tree) CreateMembershipProof(index uint64) ([]byte, error) {
	conf := tree.config

	if conf.allLeavesNum <= index {
		return nil, ErrLeafIndexOutOfRange
	}

	buf := bytes.NewBuffer(nil)

	for d := conf.depth; d > 0; d-- {
		var siblingIndex uint64
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

func (tree *Tree) VerifyMembershipProof(index uint64, proof []byte) (bool, error) {
	conf := tree.config

	if conf.allLeavesNum <= index {
		return false, ErrLeafIndexOutOfRange
	}

	if uint64(len(proof)) != conf.depth*conf.hashSize {
		return false, nil
	}

	b := tree.levels[conf.depth][index].Bytes()
	for i := uint64(0); i < conf.depth; i++ {
		sibling := proof[i*conf.hashSize : i*conf.hashSize+conf.hashSize]

		conf.hasher.Reset()
		if index%2 == 0 {
			b = concat(b, sibling)
		} else {
			b = concat(sibling, b)
		}
		if _, err := conf.hasher.Write(b); err != nil {
			return false, err
		}
		b = conf.hasher.Sum(nil)

		index /= 2
	}

	return bytes.Equal(b, tree.Root().Bytes()), nil
}
