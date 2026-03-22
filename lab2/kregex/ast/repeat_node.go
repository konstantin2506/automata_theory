package ast

import "fmt"

type RepeatNode struct {
	child Node
	count uint32
}

func (n *RepeatNode) Children() []Node {
	return []Node{n.child}
}

func (n *RepeatNode) Type() NodeT {
	return Repeat
}

func (n *RepeatNode) String() string {
	return fmt.Sprintf("Repeat: %d", n.count)
}
