package ast

type OpenParenNode struct{}

func (n OpenParenNode) Children() []Node {
	return nil
}

func (n *OpenParenNode) Type() NodeT {
	return OpenParen
}

func (n *OpenParenNode) String() string {
	return "( - node (wtf!!!)"
}
