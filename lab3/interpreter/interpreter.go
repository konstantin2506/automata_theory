package interpreter

import (
	"errors"
	"fmt"
)

type Interpreter struct {
	root        AstNode
	globalScope Scope
	CallStack
}

func NewInterpreter(treeRoot AstNode) Interpreter {
	return Interpreter{treeRoot, NewScope(nil), NewCallStack()}
}

func interpretRec(intr *Interpreter, node AstNode) (Variable, error) {
	switch n := node.(type) {
	case *FunctionCallNode:
		res, err := intr.MakeCall(n)
		if err != nil {
			if errors.Is(err, ErrStackLimit) {
				panic(err.Error())
			}
			return nil, err
		}
		return res, err
	default:
		return n.Eval(intr.TopScope())
	}
}

func Interpret(intr *Interpreter) {
	_, err := interpretRec(intr, intr.root)
	if err != nil {
		fmt.Printf("program finished with error: %s\n", err.Error())
	}
	fmt.Printf("program finished without errors\n")
}
