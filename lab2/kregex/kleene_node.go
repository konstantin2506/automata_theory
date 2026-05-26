package kregex

type KleeneNode struct {
	child Node
}

func (n *KleeneNode) Children() []Node {
	return []Node{n.child}
}

func (n *KleeneNode) Type() NodeT {
	return Kleene
}

func (n *KleeneNode) String() string {
	return "Kleene"
}

func (n *KleeneNode) CalcNullable(specMap map[Node]*NodeSpec) bool {
	n.child.(GLNode).CalcNullable(specMap)
	return SetNullable(n, true, specMap)
}

func (n *KleeneNode) CalcFirst(specMap map[Node]*NodeSpec, charNums map[Node]int) []int {
	res := n.child.(GLNode).CalcFirst(specMap, charNums)
	specMap[n].First = res
	return res
}

func (n *KleeneNode) CalcLast(specMap map[Node]*NodeSpec, charNums map[Node]int) []int {
	res := n.child.(GLNode).CalcLast(specMap, charNums)
	specMap[n].Last = res
	return res
}
