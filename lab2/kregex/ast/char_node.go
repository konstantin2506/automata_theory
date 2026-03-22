package ast

type CharNode struct {
	char byte
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
