// Package interpreter
package interpreter

import (
	"errors"
	"fmt"
)

type Scope struct {
	variables map[string]Variable
}

var (
	ErrVarDoubleDeclaration = errors.New("double declaration of variable")
	ErrVarNotDeclared       = errors.New("double declaration of variable")
	ErrVarInvalidType       = errors.New("invalid type of variable")
)

func (scope *Scope) CheckVarExistance(varName string) (bool, error) {
	_, ok := scope.variables[varName]
	if ok {
		return true, fmt.Errorf("%w: '%s'", ErrVarDoubleDeclaration, varName)
	}
	return false, fmt.Errorf("%w: '%s'", ErrVarNotDeclared, varName)
}

func (scope *Scope) ConstructInt(varName string, value int) error {
	_, err := scope.CheckVarExistance(varName)
	if err != nil {
		return fmt.Errorf("construct Int error: %w", err)
	}
	scope.variables[varName] = NewVariableInt(value)
	return nil
}

func (scope *Scope) ConstructBool(varName string, value bool) error {
	_, err := scope.CheckVarExistance(varName)
	if err != nil {
		return fmt.Errorf("construct Bool error: %w", err)
	}
	scope.variables[varName] = NewVariableBool(value)
	return nil
}

func (scope *Scope) ConstructBoolArray(varName string, sizes []int, value bool) error {
	_, err := scope.CheckVarExistance(varName)
	if err != nil {
		return fmt.Errorf("construct BoolArray error: %w", err)
	}
	boolArray, err := NewBoolArray(sizes, value)
	if err != nil {
		return fmt.Errorf("construct BoolArray error: %w", err)
	}
	scope.variables[varName] = boolArray
	return nil
}

func (scope *Scope) ConstructIntArray(varName string, sizes []int, value int) error {
	_, err := scope.CheckVarExistance(varName)
	if err != nil {
		return fmt.Errorf("construct IntArray error: %w", err)
	}
	intArray, err := NewIntArray(sizes, value)
	if err != nil {
		return fmt.Errorf("construct IntArray error: %w", err)
	}
	scope.variables[varName] = intArray
	return nil
}

func (scope *Scope) AssignInt(varName string, value int) error {
	exists, err := scope.CheckVarExistance(varName)
	if !exists {
		return fmt.Errorf("assign int error: %w", err)
	}
	v := scope.variables[varName]
	if v.Type() != Int {
		return fmt.Errorf("assign int error: %w", ErrVarInvalidType)
	}
	v.(*Integer).Assign(value)
	return nil
}

func (scope *Scope) AssignBool(varName string, value bool) error {
	exists, err := scope.CheckVarExistance(varName)
	if !exists {
		return fmt.Errorf("assign bool error: %w", err)
	}
	v := scope.variables[varName]
	if v.Type() != Bool {
		return fmt.Errorf("assign bool error: %w", ErrVarInvalidType)
	}
	v.(*Boolean).Assign(value)
	return nil
}

func (scope *Scope) AssignBoolArray(varName string, indices []int, value bool) error {
	exists, err := scope.CheckVarExistance(varName)
	if !exists {
		return fmt.Errorf("assign BoolArray error: %w", err)
	}
	v := scope.variables[varName]
	if v.Type() != Bool {
		return fmt.Errorf("assign BoolArray error: %w", ErrVarInvalidType)
	}
	err = v.(*VarArray[bool]).Assign(indices, value)
	return err
}

func (scope *Scope) AssignIntArray(varName string, indices []int, value int) error {
	exists, err := scope.CheckVarExistance(varName)
	if !exists {
		return fmt.Errorf("assign IntArray error: %w", err)
	}
	v := scope.variables[varName]
	if v.Type() != Int {
		return fmt.Errorf("assign IntArray error: %w", ErrVarInvalidType)
	}
	err = v.(*VarArray[int]).Assign(indices, value)
	return err
}
