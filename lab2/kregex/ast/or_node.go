package ast

type OrNode struct {
	childs []Node
}

func (n *OrNode) Children() []Node {
	return n.childs
}

func (n *OrNode) Type() NodeT {
	return Or
}

func (n *OrNode) String() string {
	return "Or"
}
