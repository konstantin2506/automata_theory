package main

import (
	"fmt"

	ast "kregex/kregex/ast"
)

func main() {
	str := "a"
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

	empty := ast.Difference(dfa, dfa)

	sub := ast.Difference(dfa, empty)

	complement := ast.Complement(dfa)
	sub.SaveToDot("sub.dot")
	err = complement.SaveToDot("complement.dot")
	if err != nil {
		fmt.Println(err)
	}

	toSearch := "qqqbbbbbbccaaaaaaaac"
	found, groups := nfa.Search(toSearch)
	fmt.Printf("Str: '%s' - resNFA=%s, groups = %v\n", toSearch, found, groups)
	substr := dfa.Search(toSearch)
	fmt.Printf("Str: '%s' - resDFA=%s\n", toSearch, substr)

	rx := ast.RestoreRegex(&dfa)
	fmt.Println("restored:", rx)
	rxTree, _ := ast.BuildAst(rx)
	fromRegex := ast.Minimize(ast.NewDfa(rxTree))
	fromRegex.SaveToDot("regex.dot")
	// fromRegex.PrintDebug()
	// bad, _ := ast.BuildAst("aboa")
	// b := ast.Minimize(ast.NewDfa(bad))
	fmt.Println("Are isomorphic:", ast.AreIsomorphic(&fromRegex, &dfa))
}
