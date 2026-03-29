package ast

type Node interface {
	Children() []Node
	Type() NodeT
	String() string
	CalcNullable(specMap map[Node]*NodeSpec) bool
	CalcFirst(specMap map[Node]*NodeSpec, charNums map[Node]int) []int
	// CalcLast()
}

type (
	NodeT  int
	IntSet map[int]struct{}
)

const (
	Concat NodeT = iota
	Kleene
	NamedGroup
	Optional
	Or
	Reference
	Repeat
	Char
	OpenParen
)

type NodeSpec struct {
	IsNullable bool
	First      []int
	Last       []int
}

func NewNodeSpecMap() map[Node]*NodeSpec {
	return make(map[Node]*NodeSpec)
}

func SetNullable(n Node, value bool, specMap map[Node]*NodeSpec) bool {
	nspec := &NodeSpec{}
	nspec.IsNullable = value
	specMap[n] = nspec

	return value
}
