package nodes

type ReferenceNode struct {
	refName string
}

func (n *ReferenceNode) Children() []Node {
	return nil
}

func (n *ReferenceNode) Name() string {
	return "Ref"
}

func (n *ReferenceNode) RefName() string {
	return n.refName
}
