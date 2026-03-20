package ast

type NodeStack struct {
	data []Node
}

func (s *NodeStack) isEmpty() bool {
	return (len(s.data) == 0)
}

func (s *NodeStack) Push(node Node) {
	s.data = append(s.data, node)
}

func (s *NodeStack) Pop() bool {
	if s.isEmpty() {
		return false
	}
	s.data = s.data[:len(s.data)-1]
	return true
}

func (s *NodeStack) Top() Node {
	return s.data[len(s.data)-1]
}

func (s *NodeStack) Size() int {
	return len(s.data)
}
