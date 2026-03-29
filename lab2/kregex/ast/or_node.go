package ast

import "github.com/samber/lo"

type OrNode struct {
	childs []Node
}

func (n *OrNode) Children() []Node {
	return n.childs
}

func (n *OrNode) Type() NodeT {
	return Or
}

func (n *OrNode) String() string {
	return "Or"
}

func (n *OrNode) CalcNullable(specMap map[Node]*NodeSpec) bool {
	res := false
	for _, child := range n.childs {
		res = child.CalcNullable(specMap) || res
	}
	return SetNullable(n, res, specMap)
}

func (n *OrNode) CalcFirst(specMap map[Node]*NodeSpec, charNums map[Node]int) []int {
	res := []int{}
	for _, child := range n.childs {
		res = lo.Union(res, child.CalcFirst(specMap, charNums))
	}
	specMap[n].First = res
	return res
}
