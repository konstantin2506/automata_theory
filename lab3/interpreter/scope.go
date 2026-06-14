// Package interpreter
package interpreter

import (
	"errors"
	"fmt"
)

type Scope struct {
	parent    *Scope
	variables map[string]Variable
}

var (
	ErrVarDoubleDeclaration = errors.New("double declaration of variable")
	ErrVarNotDeclared       = errors.New("double declaration of variable")
	ErrVarInvalidType       = errors.New("invalid type of variable")
	ErrNotAVectorType       = errors.New("not a vector type")
	ErrNotAScalarType       = errors.New("not a scalar type")
)

func NewScope(parent *Scope) Scope {
	return Scope{parent, make(map[string]Variable)}
}

func (scope *Scope) FindVariableDepth(varName string) (Variable, error) {
	variable, ok := scope.variables[varName]
	if ok {
		return variable, nil
	}
	if scope.parent == nil {
		return nil, fmt.Errorf("%w: '%s'", ErrVarNotDeclared, varName)
	}
	upperVar, err := scope.parent.parent.FindVariableDepth(varName)

	return upperVar, err
}

func (scope *Scope) CheckDoubleDecl(varName string) error {
	_, ok := scope.variables[varName]
	if !ok {
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
