package ast

import (
	"github.com/samber/lo"
)

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
		res = child.(GLNode).CalcNullable(specMap) || res
	}
	return SetNullable(n, res, specMap)
}

func (n *OrNode) CalcFirst(specMap map[Node]*NodeSpec, charNums map[Node]int) []int {
	res := []int{}
	for _, child := range n.childs {
		res = lo.Union(res, child.(GLNode).CalcFirst(specMap, charNums))
	}
	specMap[n].First = res
	return res
}

func (n *OrNode) CalcLast(specMap map[Node]*NodeSpec, charNums map[Node]int) []int {
	res := []int{}
	for _, child := range n.childs {
		res = lo.Union(res, child.(GLNode).CalcLast(specMap, charNums))
	}
	specMap[n].Last = res
	return res
}
