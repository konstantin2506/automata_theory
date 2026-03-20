package ast

type ReferenceNode struct {
	refName string
}

func (n *ReferenceNode) Children() []Node {
	return nil
}

func (n *ReferenceNode) Type() NodeT {
	return Reference
}

func (n *ReferenceNode) RefName() string {
	return n.refName
}
