package ast

import (
	"github.com/samber/lo"
)

type GLNode interface {
	CalcNullable(specMap map[Node]*NodeSpec) bool
	CalcFirst(specMap map[Node]*NodeSpec, charNums map[Node]int) []int
	CalcLast(specMap map[Node]*NodeSpec, charNums map[Node]int) []int
}

type (
	FollowMap   map[int][]int
	FirstMap    map[int][]int
	LastMap     map[int][]int
	NullableMap map[int]bool
)

// compile time interface cast precheck
var (
	_ GLNode = (*CharNode)(nil)
	_ GLNode = (*ConcatNode)(nil)
	_ GLNode = (*OrNode)(nil)
	_ GLNode = (*KleeneNode)(nil)
	_ GLNode = (*OptionalNode)(nil)
	_ GLNode = (*RepeatNode)(nil)
	_ GLNode = (*NamedGroupNode)(nil)
)

func markCharsDev(tree Ast, chars []byte, charNums map[Node]int) []byte {
	// RLR
	if len(chars) == 0 {
		chars = append(chars, '\n')
	}
	root := tree.GetRoot()
	if root.Type() == Char {
		charNums[root] = len(chars)
		chars = append(chars, root.String()[0])
	}
	for _, child := range root.Children() {
		subtree := NewSubAst(child)
		chars = markCharsDev(subtree, chars, charNums)
	}
	return chars
}

func MarkChars(tree *Ast) ([]byte, map[Node]int) {
	charNums := make(map[Node]int)
	chars := markCharsDev(*tree, []byte{}, charNums)

	// finish of regex marker
	endMarker := &CharNode{char: '#', number: 0}
	charNums[endMarker] = len(chars)
	chars = append(chars, '#')

	newRoot := &ConcatNode{
		children: []Node{tree.GetRoot(), endMarker},
	}
	tree.root = newRoot

	return chars, charNums
}

func ComputeNullable(tree *Ast, specMap map[Node]*NodeSpec, charNums map[Node]int) NullableMap {
	tree.GetRoot().(GLNode).CalcNullable(specMap)

	nullable := NullableMap{}
	for node, nodeSpec := range specMap {
		if node.Type() == Char {
			pos := charNums[node]
			nullable[pos] = nodeSpec.IsNullable
		}
	}
	return nullable
}

func ComputeFirst(tree *Ast, specMap map[Node]*NodeSpec, charNums map[Node]int) FirstMap {
	tree.GetRoot().(GLNode).CalcFirst(specMap, charNums)

	first := FirstMap{}
	for node, nodeSpec := range specMap {
		if node.Type() == Char {
			pos := charNums[node]
			first[pos] = nodeSpec.First
		}
	}
	return first
}

func ComputeLast(tree *Ast, specMap map[Node]*NodeSpec, charNums map[Node]int) LastMap {
	tree.GetRoot().(GLNode).CalcLast(specMap, charNums)
	last := LastMap{}
	for node, nodeSpec := range specMap {
		if node.Type() == Char {
			pos := charNums[node]
			last[pos] = nodeSpec.Last
		}
	}
	return last
}

func ComputeFollow(specMap map[Node]*NodeSpec) FollowMap {
	follow := FollowMap{}

	for node := range specMap {
		switch node.(type) {
		case *ConcatNode:
			children := node.Children()
			for i := 0; i < len(children)-1; i++ {
				left := children[i]
				for j := i + 1; j < len(children); j++ {
					right := children[j]
					for _, pos := range specMap[left].Last {
						follow[pos] = lo.Union(follow[pos], specMap[right].First)
					}
					if !specMap[right].IsNullable {
						break
					}
				}
			}
		case *KleeneNode:
			child := node.Children()[0]

			for _, pos := range specMap[child].Last {
				follow[pos] = lo.Union(follow[pos], specMap[child].First)
			}
		}
	}
	return follow
}
