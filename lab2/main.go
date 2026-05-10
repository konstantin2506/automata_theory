package main

import (
	"fmt"

	ast "kregex/kregex/ast"
)

func main() {
	str := "((<ABC>abcde...)f...g...)"
	//	str = "(ab)...|(qw)"
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
	tree.Print(1)

	dfa := ast.NewDfa(tree)
	fmt.Println(dfa.StatesCount())
	err = dfa.SaveToDot("dfa.dot")
	dfa = ast.Minimize(dfa)
	fmt.Println(dfa.StatesCount())

	err = dfa.SaveToDot("dfa_m.dot")
	if err != nil {
		fmt.Println(err)
	}

	toSearch := "qqqabcdddddeeeefg"
	found, groups := nfa.Search(toSearch)
	fmt.Printf("Str: '%s' - resNFA=%s, groups = %v\n", toSearch, found, groups)
	substr := dfa.Search(toSearch)
	fmt.Printf("Str: '%s' - resDFA=%s\n", toSearch, substr)

}
