package interpreter

import (
	"errors"
	"fmt"
)

type Interpreter struct {
	root        AstNode
	globalScope *Scope
	CallStack
}

func NewInterpreter(treeRoot AstNode) *Interpreter {
	globalScope := NewScope(nil)

	intr := &Interpreter{treeRoot, &globalScope, NewCallStack()}
	intr.scopeStack = append(intr.scopeStack, &globalScope)
	return intr
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
		fmt.Printf("[Yucky]: program finished with error: %s\n", err.Error())
		return
	}
	fmt.Printf("[Yucky]: program finished without errors\n")
}
