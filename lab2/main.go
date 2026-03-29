package main

import (
	"fmt"

	ast "kregex/kregex/ast"
	dfa "kregex/kregex/dfa"
)

func main() {
	str := "(a|b)ba?"
	//[(), ?, |, +, ()]
	fmt.Println(str)
	tree, err := ast.BuildAst(str)
	if err != nil {
		fmt.Println(err)
		return
	}
	tree.Print(1)
	specMap := ast.NewNodeSpecMap()
	_, charNums := dfa.MarkChars(tree)

	dfa.ComputeNullable(&tree, specMap)
	dfa.ComputeFirst(&tree, specMap, charNums)
	for key, value := range specMap {
		fmt.Printf("%s : %v\n", key.String(), *value)
	}
}
