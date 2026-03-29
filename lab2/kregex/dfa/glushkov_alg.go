// Package dfa
// operations and creating dfa from ast
package dfa

import (
	ast "kregex/kregex/ast"
)

func markCharsDev(tree ast.Ast, chars []byte, charNums map[ast.Node]int) []byte {
	// RLR
	if len(chars) == 0 {
		chars = append(chars, '\n')
	}
	root := tree.GetRoot()
	if root.Type() == ast.Char {
		charNums[root] = len(chars)
		chars = append(chars, root.String()[0])
	}
	for _, child := range root.Children() {
		subtree := ast.NewSubAst(child)
		chars = markCharsDev(subtree, chars, charNums)
	}
	return chars
}

func MarkChars(tree ast.Ast) ([]byte, map[ast.Node]int) {
	charNums := make(map[ast.Node]int)
	return markCharsDev(tree, []byte{}, charNums), charNums
}

func ComputeNullable(ast *ast.Ast, specMap map[ast.Node]*ast.NodeSpec) bool {
	return ast.GetRoot().CalcNullable(specMap)
}

func ComputeFirst(ast *ast.Ast, specMap map[ast.Node]*ast.NodeSpec, charNums map[ast.Node]int) {
	ast.GetRoot().CalcFirst(specMap, charNums)
}

func ComputeLast(ast *ast.Ast) {}

func ComputeFollow(ast *ast.Ast) {}

func BuildDFA() {}
