package ast

type ConcatNode struct {
	children []Node
}

func (n *ConcatNode) Children() []Node {
	return n.children
}

func (n *ConcatNode) Type() NodeT {
	return Concat
}

func (n *ConcatNode) String() string {
	return "Concat"
}
