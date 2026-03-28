package main

import (
	"fmt"

	"kregex/kregex/ast"
)

func main() {
	str := "(<abc>)"
	//[(), ?, |, +, ()]
	fmt.Println(str)
	ast, err := ast.BuildAst(str)
	if err != nil {
		fmt.Println(err)
		return
	}
	res := ast.TraverseRLR("", 1, ' ', "")
	fmt.Println(res)
	ast.Print(1)
}
