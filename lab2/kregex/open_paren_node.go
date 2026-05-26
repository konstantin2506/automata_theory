package kregex

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
