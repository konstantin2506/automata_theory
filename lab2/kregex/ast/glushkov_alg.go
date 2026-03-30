package ast

import "github.com/samber/lo"

type GLNode interface {
	CalcNullable(specMap map[Node]*NodeSpec) bool
	CalcFirst(specMap map[Node]*NodeSpec, charNums map[Node]int) []int
	CalcLast(specMap map[Node]*NodeSpec, charNums map[Node]int) []int
}

type FollowMap map[int][]int

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

func MarkChars(tree Ast) ([]byte, map[Node]int) {
	charNums := make(map[Node]int)
	return markCharsDev(tree, []byte{}, charNums), charNums
}

func ComputeNullable(tree *Ast, specMap map[Node]*NodeSpec) bool {
	return tree.GetRoot().(GLNode).CalcNullable(specMap)
}

func ComputeFirst(tree *Ast, specMap map[Node]*NodeSpec, charNums map[Node]int) {
	tree.GetRoot().(GLNode).CalcFirst(specMap, charNums)
}

func ComputeLast(tree *Ast, specMap map[Node]*NodeSpec, charNums map[Node]int) {
	tree.GetRoot().(GLNode).CalcLast(specMap, charNums)
}

func ComputeFollow(specMap map[Node]*NodeSpec) FollowMap {
	follow := make(FollowMap)
	for node := range specMap {
		switch node.(type) {
		case *ConcatNode:
			children := node.Children()
			for i := 0; i < len(children)-1; i++ {
				left := children[i]
				right := children[i+1]
				for _, pos := range specMap[left].Last {
					follow[pos] = lo.Union(follow[pos], specMap[right].First)
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

func BuildDFA() {}
