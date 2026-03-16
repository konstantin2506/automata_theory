package nodes

type RepeatNode struct {
	child Node
	count uint32
}

func (n *RepeatNode) Children() []Node {
	return []Node{n.child}
}

func (n *RepeatNode) Name() string {
	return "Repeat"
}

func (n *RepeatNode) Count() uint32 {
	return n.count
}
