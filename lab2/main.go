package main

import (
	"fmt"

	"kregex/kregex/ast"
)

func main() {
	str := "((a(<inner>b...)?|c{2}))"
	fmt.Println(str)
	ast, err := ast.BuildAst(str)
	if err != nil {
		fmt.Println(err)
		return
	}

	ast.Print(1)
}
