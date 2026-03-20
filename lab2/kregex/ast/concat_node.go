package ast

type ConcatNode struct {
	childs []Node
}

func (n *ConcatNode) Children() []Node {
	return n.childs
}

func (n *ConcatNode) Type() NodeT {
	return Concat
}
