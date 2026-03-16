// Package kregex
// Kostya's regular expressions library
package kregex

import (
	"fmt"

	nodes "kregex/kregex/nodes"
)

type AST struct {
	root nodes.Node
}

func BuildAST(str string) (AST, error) {
	var err error
	firstStr := fmt.Sprintf("(%s)", str)
	parser := NewParser(firstStr)

	err = parser.ParseParenths()
	if err != nil {
		return AST{}, fmt.Errorf("parenths parsing failed: %w", err)
	}

	fmt.Println(parser)

	for i := range parser.openParensStack {
		first, last, err := parser.GetNextParenths(i)
		if err != nil {
			break
		}
		fmt.Printf("{%d, %d}\n", first, last)
		parser.ParseInside(first, last)
		fmt.Println()
	}

	return AST{nil}, nil
}
