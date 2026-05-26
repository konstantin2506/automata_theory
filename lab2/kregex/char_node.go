package kregex

type CharNode struct {
	char   byte
	number int
}

func (n *CharNode) Children() []Node {
	return nil
}

func (n *CharNode) Type() NodeT {
	return Char
}

func (n *CharNode) String() string {
	return string(n.char)
}

func (n *CharNode) CalcNullable(specMap map[Node]*NodeSpec) bool {
	return SetNullable(n, false, specMap)
}

func (n *CharNode) CalcFirst(specMap map[Node]*NodeSpec, charNums map[Node]int) []int {
	res := []int{}
	res = append(res, charNums[n])
	specMap[n].First = res
	return res
}

func (n *CharNode) CalcLast(specMap map[Node]*NodeSpec, charNums map[Node]int) []int {
	res := []int{}
	res = append(res, charNums[n])
	specMap[n].Last = res
	return res
}
