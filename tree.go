package merkle

import (
	"bytes"
	"errors"
	"hash"
	"math/big"
)

const (
	DepthMax = 32
)

var (
	ErrTooLargeDepth     = errors.New("too large tree depth")
	ErrTooManyLeaves     = errors.New("too many leaves")
	ErrTooLargeLeafIndex = errors.New("too large leaf index")
	ErrInvalidProofSize  = errors.New("invalid proof size")
)

type Tree struct {
	hasher       hash.Hash
	hashSize     uint64
	depth        uint64
	allLeavesNum uint64
	allNodesNum  uint64
	nodes        []*Node
	levels       [][]*Node
}

func NewTree(hasher hash.Hash, depth uint64, leaves [][]byte) (*Tree, error) {
	if depth > DepthMax {
		return nil, ErrTooLargeDepth
	}

	allLeavesNum := new(big.Int).Lsh(big.NewInt(2), uint(depth-1)).Uint64()
	if uint64(len(leaves)) > allLeavesNum {
		return nil, ErrTooManyLeaves
	}

	allNodesNum := uint64(0)
	for i := allLeavesNum; i >= 1; i /= 2 {
		allNodesNum += i
	}

	tree := &Tree{
		hasher:       hasher,
		hashSize:     uint64(hasher.Size()),
		depth:        depth,
		allLeavesNum: allLeavesNum,
		allNodesNum:  allNodesNum,
		nodes:        make([]*Node, allNodesNum),
		levels:       make([][]*Node, depth+1),
	}

	if err := tree.buildBase(leaves); err != nil {
		return nil, err
	}
	if err := tree.build(); err != nil {
		return nil, err
	}

	return tree, nil
}

func (tree *Tree) hash(b []byte) ([]byte, error) {
	tree.hasher.Reset()
	if _, err := tree.hasher.Write(b); err != nil {
		return nil, err
	}
	return tree.hasher.Sum(nil), nil
}

func (tree *Tree) pairHash(b1, b2 []byte) ([]byte, error) {
	tree.hasher.Reset()
	if _, err := tree.hasher.Write(b1); err != nil {
		return nil, err
	}
	if _, err := tree.hasher.Write(b2); err != nil {
		return nil, err
	}
	return tree.hasher.Sum(nil), nil
}

func (tree *Tree) buildBase(leaves [][]byte) error {
	leavesNum := uint64(len(leaves))

	tree.levels[tree.depth] = make([]*Node, tree.allLeavesNum)

	for i := uint64(0); i < leavesNum; i++ {
		leafHash, err := tree.hash(leaves[i])
		if err != nil {
			return err
		}

		node := newNode(leafHash, nil, nil)
		tree.nodes[i] = node
		tree.levels[tree.depth][i] = node
	}

	emptyLeafHash, err := tree.hash(make([]byte, tree.hashSize, tree.hashSize))
	if err != nil {
		return err
	}

	for i := leavesNum; i < tree.allLeavesNum; i++ {
		node := newNode(emptyLeafHash, nil, nil)
		tree.nodes[i] = node
		tree.levels[tree.depth][i] = node
	}

	return nil
}

func (tree *Tree) build() error {
	for d := tree.depth; d > 0; d-- {
		level := tree.levels[d]

		nextDepth := d - 1
		tree.levels[nextDepth] = make([]*Node, len(level)/2)

		for i := 0; i < len(level); i += 2 {
			left := level[i]
			right := level[i+1]

			parentHash, err := tree.pairHash(left.Bytes(), right.Bytes())
			if err != nil {
				return err
			}

			tree.levels[nextDepth][i/2] = newNode(parentHash, left, right)
		}

		tree.nodes = append(tree.nodes, tree.levels[nextDepth]...)
	}

	return nil
}

func (tree *Tree) Root() *Node {
	return tree.levels[0][0]
}

func (tree *Tree) CreateMembershipProof(index uint64) ([]byte, error) {
	if index >= tree.allLeavesNum {
		return nil, ErrTooLargeLeafIndex
	}

	buf := bytes.NewBuffer(nil)

	for d := tree.depth; d > 0; d-- {
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
	if index >= tree.allLeavesNum {
		return false, ErrTooLargeLeafIndex
	}
	if uint64(len(proof)) != tree.depth*tree.hashSize {
		return false, ErrInvalidProofSize
	}

	proofIndex := uint64(0)

	b := tree.levels[tree.depth][index].Bytes()

	for d := tree.depth; d > 0; d-- {
		sibling := proof[proofIndex : proofIndex+tree.hashSize]

		var err error
		if index%2 == 0 {
			b, err = tree.pairHash(b, sibling)
		} else {
			b, err = tree.pairHash(sibling, b)
		}
		if err != nil {
			return false, err
		}

		proofIndex += tree.hashSize
		index /= 2
	}

	return bytes.Equal(b, tree.Root().Bytes()), nil
}
