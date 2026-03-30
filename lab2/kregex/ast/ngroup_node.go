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

func (n *NamedGroupNode) CalcNullable(specMap map[Node]*NodeSpec) bool {
	res := n.child.(GLNode).CalcNullable(specMap)
	return SetNullable(n, res, specMap)
}

func (n *NamedGroupNode) CalcFirst(specMap map[Node]*NodeSpec, charNums map[Node]int) []int {
	res := n.child.(GLNode).CalcFirst(specMap, charNums)
	specMap[n].First = res
	return res
}

func (n *NamedGroupNode) CalcLast(specMap map[Node]*NodeSpec, charNums map[Node]int) []int {
	res := n.child.(GLNode).CalcLast(specMap, charNums)
	specMap[n].First = res
	return res
}
