package main

import (
	"fmt"

	"kregex/kregex/ast"
)

func main() {
	str := "((a|b)|(<Q>q{10})(<C>c...))"
	//[(), ?, |, +, ()]
	fmt.Println(str)
	ast, err := ast.BuildAst(str)
	if err != nil {
		fmt.Println(err)
		return
	}

	ast.Print(1)
}
