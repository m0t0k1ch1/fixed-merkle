package merkle

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"math"
)

const (
	DepthMin    = 1
	DepthMax    = 16
	LeafSizeMin = 1  // bytes
	LeafSizeMax = 64 // bytes
)

var (
	ErrTooSmallDepth    = fmt.Errorf("depth must be %d or more", DepthMin)
	ErrTooLargeDepth    = fmt.Errorf("depth must be %d or less", DepthMax)
	ErrTooSmallLeafSize = fmt.Errorf("leaf size must be %d or more", LeafSizeMin)
	ErrTooLargeLeafSize = fmt.Errorf("leaf size must be %d or less", LeafSizeMax)
	ErrTooManyLeaves    = errors.New("number of leaves exceeds upper limit")
)

type TreeBuilder struct {
	hasher    hash.Hash
	depth     int
	leafSize  int
	leavesNum int
	nodesNum  int
}

var DefaultTreeBuilder = &TreeBuilder{
	hasher:    sha256.New(),
	depth:     16,
	leafSize:  32,
	leavesNum: 65536,
	nodesNum:  131071,
}

func NewTreeBuilder(hasher hash.Hash, depth int, leafSize int) (*TreeBuilder, error) {
	if depth < DepthMin {
		return nil, ErrTooSmallDepth
	}
	if depth > DepthMax {
		return nil, ErrTooLargeDepth
	}
	if leafSize < LeafSizeMin {
		return nil, ErrTooSmallLeafSize
	}
	if leafSize > LeafSizeMax {
		return nil, ErrTooLargeLeafSize
	}

	leavesNum := int(math.Exp2(float64(depth)))

	nodesNum := 0
	for i := leavesNum; i >= 1; i /= 2 {
		nodesNum += i
	}

	return &TreeBuilder{
		hasher:    hasher,
		depth:     depth,
		leafSize:  leafSize,
		leavesNum: leavesNum,
		nodesNum:  nodesNum,
	}, nil
}

func (builder *TreeBuilder) Build(leaves [][]byte, hashed bool) (*Tree, error) {
	leavesNum := len(leaves)
	if leavesNum > builder.leavesNum {
		return nil, ErrTooManyLeaves
	}

	if !hashed {
		for i, leaf := range leaves {
			builder.hasher.Reset()
			if _, err := builder.hasher.Write(leaf); err != nil {
				return nil, err
			}
			leaves[i] = builder.hasher.Sum(nil)
		}
	}

	nodes := make([]*Node, builder.nodesNum)
	bottomLevel := make([]*Node, builder.leavesNum)

	for i := 0; i < leavesNum; i++ {
		node := NewNode(leaves[i], nil, nil)
		nodes[i] = node
		bottomLevel[i] = node
	}

	for i := leavesNum; i < builder.leavesNum; i++ {
		node := NewNode(make([]byte, builder.leafSize, builder.leafSize), nil, nil)
		nodes[i] = node
		bottomLevel[i] = node
	}

	levels := make([][]*Node, builder.depth+1)
	levels[builder.depth] = bottomLevel

	tree := &Tree{
		Nodes:  nodes,
		Levels: levels,
	}

	for depth := builder.depth; depth > 0; depth-- {
		builder.buildNextLevel(tree, depth)
	}

	return tree, nil
}

func (builder *TreeBuilder) buildNextLevel(tree *Tree, currentDepth int) error {
	currentLevel := tree.Levels[currentDepth]

	nextDepth := currentDepth - 1
	tree.Levels[nextDepth] = make([]*Node, len(currentLevel)/2)

	for i := 0; i < len(currentLevel); i += 2 {
		left := currentLevel[i]
		right := currentLevel[i+1]

		builder.hasher.Reset()

		b := make([]byte, len(left.HashBytes)+len(right.HashBytes))
		copy(b[0:len(left.HashBytes)], left.HashBytes[:])
		copy(b[len(left.HashBytes):len(b)], right.HashBytes[:])

		if _, err := builder.hasher.Write(b); err != nil {
			return err
		}
		hashBytes := builder.hasher.Sum(nil)

		tree.Levels[nextDepth][i/2] = NewNode(hashBytes, left, right)
	}

	tree.Nodes = append(tree.Nodes, tree.Levels[nextDepth]...)

	return nil
}
