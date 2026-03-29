package ast

type OpenParenNode struct{}

func (n OpenParenNode) Children() []Node {
	return nil
}

func (n *OpenParenNode) Type() NodeT {
	return OpenParen
}

func (n *OpenParenNode) String() string {
	panic("String() : (-node (wtf!!!)")
}

func (n *OpenParenNode) CalcNullable(specMap map[Node]*NodeSpec) bool {
	panic("Open Paren in 'correct' ast")
}

func (n *OpenParenNode) CalcFirst(specMap map[Node]*NodeSpec, charNums map[Node]int) []int {
	panic("Open Paren in 'correct ast")
}
