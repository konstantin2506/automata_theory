package main

import (
	"fmt"

	"kregex/kregex/ast"
)

func main() {
	str := "(abc(qx)?(please)?)"
	fmt.Println(str)
	ast, err := ast.BuildAst(str)
	if err != nil {
		fmt.Println(err)
		return
	}

	ast.Print(1)
}
