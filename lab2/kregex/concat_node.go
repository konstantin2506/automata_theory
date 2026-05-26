package kregex

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
		glNode := child.(GLNode)
		res = glNode.CalcNullable(specMap) && res
	}
	return SetNullable(n, res, specMap)
}

func (n *ConcatNode) CalcFirst(specMap map[Node]*NodeSpec, charNums map[Node]int) []int {
	res := []int{}
	i := 0
	for j, child := range n.children {
		res = lo.Union(res, child.(GLNode).CalcFirst(specMap, charNums))
		if !specMap[child].IsNullable {
			i = j + 1
			break
		}
	}
	for ; i < len(n.children); i++ {
		n.children[i].(GLNode).CalcFirst(specMap, charNums)
	}
	specMap[n].First = res
	return res
}

func (n *ConcatNode) CalcLast(specMap map[Node]*NodeSpec, charNums map[Node]int) []int {
	res := []int{}
	j := 0
	for i := len(n.children) - 1; i >= 0; i-- {
		child := n.children[i]
		res = lo.Union(res, child.(GLNode).CalcLast(specMap, charNums))
		if !specMap[child].IsNullable {
			j = i - 1
			break
		}
	}
	for ; j >= 0; j-- {
		n.children[j].(GLNode).CalcLast(specMap, charNums)
	}
	specMap[n].Last = res
	return res
}
