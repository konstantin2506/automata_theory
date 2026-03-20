package main

import (
	"fmt"

	"kregex/kregex/ast"
)

func main() {
	str := "((abc(q+))df(ijk)*)"
	fmt.Println(str)
	ast.BuildAst(str)
}
