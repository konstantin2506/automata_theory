package main

import (
	"fmt"

	ast "kregex/kregex/ast"
)

func main() {
	str := "(<A>a...)(<B>b?)(<C>ccc)"
	//[(), ?, |, +, ()]
	fmt.Println(str)
	tree, err := ast.BuildAst(str)
	if err != nil {
		fmt.Println(err)
		return
	}
	nfa := ast.BuildFromAst(tree.GetRoot())
	nfa.SaveToDot("nfa.dot")
	tree.SaveToDot("ast.dot")
	specMap := ast.NewNodeSpecMap()
	chars, charNums := ast.MarkChars(&tree)
	tree.Print(1)

	ast.ComputeNullable(&tree, specMap, charNums)
	first := ast.ComputeFirst(&tree, specMap, charNums)
	last := ast.ComputeLast(&tree, specMap, charNums)

	follow := ast.ComputeFollow(specMap)

	for node, spec := range specMap {
		fmt.Println(node.String(), spec.First, spec.Last)
	}
	fmt.Println()
	for x, follows := range follow {
		fmt.Println(x, follows)
	}
	dfa := ast.BuildDFA(tree.GetRoot(), follow, first, last,
		chars, charNums, specMap)
	fmt.Println(dfa.String())
	err = dfa.SaveGraphViz("dfa.dot")
	if err != nil {
		fmt.Println(err)
	}

	toSearch := "aaaaaaaaaaaaaaaaaccc"
	ok, groups := nfa.Search(toSearch)
	fmt.Printf("Str: '%s' - res=%t, groups = %v\n", toSearch, ok, groups)

}
