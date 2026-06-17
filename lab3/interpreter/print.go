package interpreter

import "fmt"

type Printer interface {
	Print()
}

type PrintNode struct {
	value AstNode
}

func NewPrintNode(value AstNode) *PrintNode {
	return &PrintNode{value}
}

func (node *PrintNode) Eval(scope *Scope) (Variable, error) {
	v, err := node.value.Eval(scope)
	if err != nil {
		return nil, err
	}
	v.(Printer).Print()
	fmt.Printf("\n")
	return nil, nil
}
