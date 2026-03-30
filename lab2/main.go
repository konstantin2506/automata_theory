package main

import (
	"fmt"

	ast "kregex/kregex/ast"
)

func main() {
	str := "(a|b)(c|d)"
	//[(), ?, |, +, ()]
	fmt.Println(str)
	tree, err := ast.BuildAst(str)
	if err != nil {
		fmt.Println(err)
		return
	}
	tree.Print(1)
	specMap := ast.NewNodeSpecMap()
	_, charNums := ast.MarkChars(tree)

	ast.ComputeNullable(&tree, specMap)
	ast.ComputeFirst(&tree, specMap, charNums)
	ast.ComputeLast(&tree, specMap, charNums)
	/*for key, value := range specMap {
		fmt.Printf("%s : %v\n", key.String(), *value)
	}*/

	followMap := ast.ComputeFollow(specMap)
	for key, value := range followMap {
		fmt.Printf("Follow[%d] = %v\n", key, value)
	}
}
