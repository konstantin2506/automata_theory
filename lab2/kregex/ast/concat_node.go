package ast

import (
	lo "github.com/samber/lo"
)

type ConcatNode struct {
	children []Node
}

func (n *ConcatNode) Children() []Node {
	return n.children
}

func (n *ConcatNode) Type() NodeT {
	return Concat
}

func (n *ConcatNode) String() string {
	return "Concat"
}

func (n *ConcatNode) CalcNullable(specMap map[Node]*NodeSpec) bool {
	res := true
	for _, child := range n.children {
		res = child.CalcNullable(specMap) && res
	}
	return SetNullable(n, res, specMap)
}

func (n *ConcatNode) CalcFirst(specMap map[Node]*NodeSpec, charNums map[Node]int) []int {
	res := []int{}
	for _, child := range n.children {
		res = lo.Union(res, child.CalcFirst(specMap, charNums))
		if !specMap[child].IsNullable {
			break
		}
	}
	specMap[n].First = res
	return res
}
