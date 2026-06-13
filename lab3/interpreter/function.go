package interpreter

import (
	"errors"
	"fmt"
)

const (
	stackLimit = 1000
)

var (
	ErrStackLimit               = errors.New("stack limit got")
	ErrFunctionParamTypesDiffer = errors.New("param types of function decl and function call differ")
	ErrFunctionParamCountDiffer = errors.New("param count of function decl and function call differ")
)

type AstNode interface {
	Eval(scope *Scope) (Variable, error)
}

type FunctionDeclNode struct {
	child      AstNode
	name       string
	params     []Variable
	paramNames []string
}

type FunctionCallNode struct {
	scope    *Scope
	function *FunctionDeclNode
	params   []AstNode
}

func checkTypesOfParams(got []Variable, want []Variable) error {
	if len(got) != len(want) {
		return fmt.Errorf("%w: want=%d, got=%d", ErrFunctionParamCountDiffer, len(want), len(got))
	}
	for i := range got {
		if want[i].Type() != got[i].Type() {
			return fmt.Errorf("%w: param[%d]: want=%d, got=%d", ErrFunctionParamTypesDiffer, i, want[i].Type(), got[i].Type())
		}
		wantVec, okWant := want[i].(Vector)
		gotVec, _ := got[i].(Vector)
		if okWant && (wantVec.InnerType() != gotVec.InnerType()) {
			return fmt.Errorf("%w (inner): param[%d]: want=%d, got=%d", ErrFunctionParamTypesDiffer, i, wantVec.InnerType(), gotVec.InnerType())
		}
	}
	return nil
}

func NewFunctionCallNode(function *FunctionDeclNode, paramNodes []AstNode) FunctionCallNode {
	//TODO wtf with params
	//
	return FunctionCallNode{nil, function, paramNodes}
}

func (fn *FunctionCallNode) Eval(scope *Scope) (Variable, error) {
	params := make([]Variable, len(fn.params))

	for i, node := range fn.params {
		param, err := node.Eval(scope)
		if err != nil {
			return nil, err
		}
		params[i] = param
	}
	err := checkTypesOfParams(params, fn.function.params)
	if err != nil {
		return nil, err
	}

	for i, param := range params {
		err := scope.ConstructCopy(fn.function.paramNames[i], param)
		if err != nil {
			return nil, err
		}
	}

	result, err := fn.function.child.Eval(fn.scope)
	return result, err
}

func (fn *FunctionCallNode) SetScope(scope *Scope) {
	fn.scope = scope
}

type CallStack struct {
	data       []*FunctionCallNode
	scopeStack []*Scope
	result     Variable
}

func NewCallStack() CallStack {
	return CallStack{make([]*FunctionCallNode, 0, stackLimit), make([]*Scope, 0, stackLimit), nil}
}

func (cs *CallStack) MakeCall(call *FunctionCallNode) (Variable, error) {
	newScope := NewScope(nil)
	call.SetScope(&newScope)
	cs.data = append(cs.data, call)
	cs.scopeStack = append(cs.scopeStack, &newScope)
	if len(cs.data) == stackLimit {
		return nil, fmt.Errorf("%w : %d", ErrStackLimit, stackLimit)
	}
	result, err := cs.data[len(cs.data)-1].Eval(&newScope)
	cs.scopeStack = cs.scopeStack[:len(cs.scopeStack)-1]
	cs.data = cs.data[:len(cs.data)-1]
	if err != nil {
		err = fmt.Errorf("in function %s(): %w", call.function.name, err)
	}
	return result, err
}

func (cs *CallStack) TopScope() *Scope {
	return cs.scopeStack[len(cs.scopeStack)-1]
}
