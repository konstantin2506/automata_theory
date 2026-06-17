// Package interpreter
package interpreter

import (
	"errors"
	"fmt"
)

type Game interface {
	Move(n int) error
	Surroundings() [3][][2]bool
	Rotate(n int)
}

type GlobalScope struct {
	scope     *Scope
	functions map[string]*FunctionDeclNode
	game      Game
}

func NewGlobalScope(game Game) *GlobalScope {
	scope := NewScope(nil, nil)
	return &GlobalScope{&scope, make(map[string]*FunctionDeclNode), game}
}

type Scope struct {
	parent      *Scope
	variables   map[string]Variable
	globalScope *GlobalScope
}

var (
	ErrVarDoubleDeclaration     = errors.New("double declaration of variable")
	ErrVarNotDeclared           = errors.New("variable not declared")
	ErrVarInvalidType           = errors.New("invalid type of variable")
	ErrNotAVectorType           = errors.New("not a vector type")
	ErrNotAScalarType           = errors.New("not a scalar type")
	ErrFuncDeclNotInGlobalScope = errors.New("function declaration not in gloabal scope")
	ErrFuncDoubleDecl           = errors.New("double declaration of function")
	ErrFuncNotDecl              = errors.New("function not declared")
)

func NewScope(parent *Scope, globalScope *GlobalScope) Scope {
	return Scope{parent, make(map[string]Variable), globalScope}
}

func (scope *Scope) DeclFunction(name string, node *FunctionDeclNode) error {
	if scope != scope.globalScope.scope {
		return fmt.Errorf("%w: name='%s'", ErrFuncDeclNotInGlobalScope, name)
	}
	_, exists := scope.globalScope.functions[name]
	if exists {
		return fmt.Errorf("%w: name='%s'", ErrFuncDoubleDecl, name)
	}
	scope.globalScope.functions[name] = node
	return nil
}

func (scope *Scope) FindFunctionDecl(name string) (*FunctionDeclNode, error) {
	fn, exists := scope.globalScope.functions[name]
	if !exists {
		return nil, fmt.Errorf("%w: name='%s'", ErrFuncNotDecl, name)
	}
	return fn, nil
}

func (scope *Scope) FindVariableDepth(varName string) (Variable, error) {
	variable, ok := scope.variables[varName]
	if ok {
		return variable, nil
	}
	if scope.parent == nil {
		return nil, fmt.Errorf("%w: '%s'", ErrVarNotDeclared, varName)
	}
	upperVar, err := scope.parent.FindVariableDepth(varName)

	return upperVar, err
}

func (scope *Scope) CheckDoubleDecl(varName string) error {
	_, ok := scope.variables[varName]
	if ok {
		return fmt.Errorf("%w: '%s'", ErrVarDoubleDeclaration, varName)
	}

	return nil
}

func (scope *Scope) ConstructCopy(varName string, value Variable) error {
	err := scope.CheckDoubleDecl(varName)
	if err != nil {
		return fmt.Errorf("construct copy (%s) error: %w", varName, err)
	}
	scope.variables[varName] = value.Copy()
	return nil
}

func (scope *Scope) ConstructScalar(varName string, value Variable) error {
	err := scope.CheckDoubleDecl(varName)
	if err != nil {
		return fmt.Errorf("construct %s (%s) error: %w", TypeName(value.Type()), varName, err)
	}
	scope.variables[varName] = value.Copy()
	return nil
}

func (scope *Scope) ConstructArray(varName string, sizes []int, value Variable) error {
	err := scope.CheckDoubleDecl(varName)
	if err != nil {
		return fmt.Errorf("construct Array (%s) error: %w", varName, err)
	}
	array, err := NewArray(sizes, value)
	if err != nil {
		return fmt.Errorf("construct Array (%s) error: %w", varName, err)
	}
	scope.variables[varName] = array
	return nil
}

func (scope *Scope) AssignScalar(varName string, value Variable) error {
	v, err := scope.FindVariableDepth(varName)
	if err != nil {
		return fmt.Errorf("assign scalar (%s) error: %w", varName, err)
	}
	if v.Type() != value.Type() {
		return fmt.Errorf("assign scalar (%s) error: %w", varName, ErrVarInvalidType)
	}
	scalar, ok := v.(Scalar)
	if !ok {
		return fmt.Errorf("assign Array (%s) error: %w", varName, ErrNotAScalarType)
	}

	scalar.Assign(value)
	return nil
}

func (scope *Scope) AssignVector(varName string, indices []int, value Variable) error {
	v, err := scope.FindVariableDepth(varName)
	if err != nil {
		return fmt.Errorf("assign Array (%s) error: %w", varName, err)
	}
	vec, ok := v.(Vector)
	if !ok {
		return fmt.Errorf("assign Array (%s) error: %w", varName, ErrNotAVectorType)
	}
	if !CmpTypeWithInner(vec, value) {
		return fmt.Errorf("assign Array (%s) error: %w", varName, ErrVarInvalidType)
	}

	err = vec.Assign(indices, value)
	return err
}
