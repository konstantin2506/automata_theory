package nodes

type ConcatNode struct {
	left  Node
	right Node
}

func (n *ConcatNode) Children() []Node {
	return []Node{n.left, n.right}
}
