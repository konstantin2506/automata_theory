package ast

import "fmt"

// R"(<varible_name>[a-zA-Z])" видимо

type NamedGroupNode struct {
	child Node
	name  string
}

func (n *NamedGroupNode) Children() []Node {
	return []Node{n.child}
}

func (n *NamedGroupNode) Type() NodeT {
	return NamedGroup
}

func (n *NamedGroupNode) String() string {
	return fmt.Sprintf("Named_group[%s]", n.name)
}
