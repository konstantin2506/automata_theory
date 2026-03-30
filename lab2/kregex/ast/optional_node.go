package ast

type OptionalNode struct {
	child Node
}

func (n *OptionalNode) Children() []Node {
	return []Node{n.child}
}

func (n *OptionalNode) Type() NodeT {
	return Optional
}

func (n *OptionalNode) String() string {
	return "Optional"
}

func (n *OptionalNode) CalcNullable(specMap map[Node]*NodeSpec) bool {
	n.child.(GLNode).CalcNullable(specMap)
	return SetNullable(n, true, specMap)
}

func (n *OptionalNode) CalcFirst(specMap map[Node]*NodeSpec, charNums map[Node]int) []int {
	res := n.child.(GLNode).CalcFirst(specMap, charNums)
	specMap[n].First = res
	return res
}

func (n *OptionalNode) CalcLast(specMap map[Node]*NodeSpec, charNums map[Node]int) []int {
	res := n.child.(GLNode).CalcLast(specMap, charNums)
	specMap[n].Last = res
	return res
}
