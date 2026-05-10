package main

import (
	"fmt"

	ast "kregex/kregex/ast"
)

func main() {
	str := "(a|b|c|d)...(ab)|(cd)|(ef){2}"
	str = "(ab)...|(qw)"
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
	err = dfa.SaveToDot("dfa.dot")
	if err != nil {
		fmt.Println(err)
	}

	toSearch := "baaa"
	ok, groups := nfa.Search(toSearch)
	fmt.Printf("Str: '%s' - resNFA=%t, groups = %v\n", toSearch, ok, groups)
	substr := dfa.Search(toSearch)
	fmt.Printf("Str: '%s' - resDFA=%s\n", toSearch, substr)

}
