package nodes

type OptionalNode struct {
	child Node
}

func (n *OptionalNode) Children() []Node {
	return []Node{n.child}
}

func (n *OptionalNode) Name() string {
	return "Optional"
}
