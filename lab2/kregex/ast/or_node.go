package ast

type OrNode struct {
	left  Node
	right Node
}

func (n *OrNode) Children() []Node {
	return []Node{n.left, n.right}
}

func (n *OrNode) Type() NodeT {
	return Or
}

func (n *OrNode) String() string {
	return "Or"
}
