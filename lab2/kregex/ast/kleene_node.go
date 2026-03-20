package ast

type KleeneNode struct {
	child Node
}

func (n *KleeneNode) Children() []Node {
	return []Node{n.child}
}

func (n *KleeneNode) Type() NodeT {
	return Kleene
}
