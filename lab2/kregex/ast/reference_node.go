package ast

import "fmt"

type ReferenceNode struct {
	refName string
}

func (n *ReferenceNode) Children() []Node {
	return nil
}

func (n *ReferenceNode) Type() NodeT {
	return Reference
}

func (n *ReferenceNode) String() string {
	return fmt.Sprintf("ref: %s", n.refName)
}
