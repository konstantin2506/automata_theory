package ast

type OptionalNode struct {
	child Node
}

func (n *OptionalNode) Children() []Node {
	return []Node{n.child}
}

func (n *OptionalNode) Type() NodeT {
	return Optional
}

func (n OptionalNode) String() string {
	return "Optional"
}
