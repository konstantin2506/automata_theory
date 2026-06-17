package interpreter

import (
	"errors"
	"fmt"
)

var (
	ErrStackLimit                = errors.New("stack limit got")
	ErrFunctionParamTypesDiffer  = errors.New("param types of function decl and function call differ")
	ErrFunctionParamCountDiffer  = errors.New("param count of function decl and function call differ")
	ErrFunctionResultTypesDiffer = errors.New("result types of function decl and function call differ")
	ErrFunctionWithoutReturn     = errors.New("function without return statement")
)

type AstNode interface {
	Eval(scope *Scope) (Variable, error)
}

type FunctionDeclNode struct {
	child      AstNode
	name       string
	params     []Variable
	paramNames []string
	result     Variable
}

func NewFunctionDeclNode(name string, params []Variable, paramNames []string, result Variable, child AstNode) *FunctionDeclNode {
	return &FunctionDeclNode{child, name, params, paramNames, result}
}

func (node *FunctionDeclNode) Eval(scope *Scope) (Variable, error) {
	if node.name == "pathfinder" {
		res, err := node.child.Eval(scope)
		if err != nil {
			return NewVariableInt(2), err
		}
		if res == nil {
			return NewVariableInt(1), fmt.Errorf("%w: 'pathfinder'", ErrFunctionWithoutReturn)
		}
		return NewVariableInt(0), nil
	}
	err := scope.DeclFunction(node.name, node)
	return nil, err
}

type FunctionCallNode struct {
	functionName string
	function     *FunctionDeclNode
	params       []AstNode
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

func NewFunctionCallNode(name string, paramNodes []AstNode) *FunctionCallNode {
	return &FunctionCallNode{name, nil, paramNodes}
}

func (fn *FunctionCallNode) Eval(scope *Scope) (Variable, error) {
	decl, err := scope.FindFunctionDecl(fn.functionName)
	if err != nil {
		return nil, err
	}
	fn.function = decl
	params := make([]Variable, len(fn.params))

	funcScope := NewScope(nil, scope.globalScope)
	for i, node := range fn.params {
		param, err := node.Eval(scope)
		if err != nil {
			return nil, err
		}
		params[i] = param
	}
	err = checkTypesOfParams(params, fn.function.params)
	if err != nil {
		return nil, err
	}

	for i, param := range params {
		err := funcScope.ConstructCopy(fn.function.paramNames[i], param)
		if err != nil {
			return nil, err
		}
	}

	result, err := fn.function.child.Eval(&funcScope)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("%w: name='%s'", ErrFunctionWithoutReturn, fn.functionName)
	}
	if fn.function.result != nil {
		if result.Type() != fn.function.result.Type() {
			return nil, fmt.Errorf("%w: want=%s, got=%s", ErrFunctionResultTypesDiffer, TypeName(fn.function.result.Type()), TypeName(result.Type()))
		}
		if vres, ok := fn.function.result.(Vector); ok {
			if vres.InnerType() != result.(Vector).InnerType() {
				return nil, fmt.Errorf("%w: want=%s, got=%s", ErrFunctionResultTypesDiffer, TypeName(vres.InnerType()), TypeName(result.(Vector).InnerType()))
			}
		}
	}

	return result, nil
}

/*
type CallStack struct {
	data       []*FunctionCallNode
	scopeStack []*Scope
	result     Variable
}

func NewCallStack() CallStack {
	return CallStack{make([]*FunctionCallNode, 0, stackLimit), make([]*Scope, 0, stackLimit), nil}
}

func (cs *CallStack) MakeCall(call *FunctionCallNode) (Variable, error) {
	newScope := NewScope(nil, cs.TopScope().globalScope)
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
}*/
