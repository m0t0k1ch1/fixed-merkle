package merkle

type Tree struct {
	Nodes  []*Node
	Levels [][]*Node
}

func (tree *Tree) Root() *Node {
	return tree.Levels[0][0]
}
