package ast

type ByteStack struct {
	data []byte
}

func (s *ByteStack) isEmpty() bool {
	return (len(s.data) == 0)
}

func (s *ByteStack) Push(elem byte) {
	s.data = append(s.data, elem)
}

func (s *ByteStack) Pop() bool {
	if s.isEmpty() {
		return false
	}
	s.data = s.data[:len(s.data)-1]
	return true
}

func (s *ByteStack) Top() byte {
	return s.data[len(s.data)-1]
}

func (s *ByteStack) Size() int {
	return len(s.data)
}
