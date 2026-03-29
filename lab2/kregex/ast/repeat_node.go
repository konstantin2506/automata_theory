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
	return fmt.Sprintf("Repeat:%d", n.count)
}

func (n *RepeatNode) CalcNullable(specMap map[Node]*NodeSpec) bool {
	if n.count == 0 {
		return SetNullable(n, true, specMap)
	}
	res := n.child.CalcNullable(specMap)
	return SetNullable(n, res, specMap)
}

func (n *RepeatNode) CalcFirst(specMap map[Node]*NodeSpec, charNums map[Node]int) []int {
	if n.count == 0 {
		return []int{}
	}
	res := n.child.CalcFirst(specMap, charNums)
	specMap[n].First = res
	return res
}
