package nodes

// R"(<varible_name>[a-zA-Z])" видимо

type NamedGroupNode struct {
	child Node
	name  string
}

func (n *NamedGroupNode) Children() []Node {
	return []Node{n.child}
}

func (n *NamedGroupNode) Name() string {
	return "ngroup"
}

func (n *NamedGroupNode) GroupName() string {
	return n.name
}
