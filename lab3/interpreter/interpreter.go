package interpreter

import (
	"fmt"
)

type Interpreter struct {
	root        AstNode
	globalScope *GlobalScope
	functions   map[string]*FunctionDeclNode
	good        int
}

func NewInterpreter(treeRoot AstNode, good int) *Interpreter {
	globalScope := NewGlobalScope(nil)
	globalScope.scope.globalScope = globalScope
	intr := &Interpreter{treeRoot, globalScope, make(map[string]*FunctionDeclNode), good}
	return intr
}

func (intr *Interpreter) SetGame(game Game) {
	intr.globalScope.game = game
}

func Interpret(intr *Interpreter, game Game) {
	intr.SetGame(game)
	_, err := intr.root.Eval(intr.globalScope.scope)
	if err != nil {
		fmt.Printf("[Yucky]: program finished with error: %s\n", err.Error())
		return
	}
	fmt.Printf("[Yucky]: program finished without errors\n")
}
