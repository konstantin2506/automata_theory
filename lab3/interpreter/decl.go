package interpreter

import (
	"errors"
	"fmt"
)

var (
	ErrScalarDeclTypesDiffer = errors.New("scalar decl types differ")
	ErrSizeNotInt            = errors.New("size of dimension is not integer")
)

type ScalarDeclNode struct {
	name   string
	value  AstNode
	innerT VarT
}

func NewScalarDeclNode(name string, value AstNode, innerT VarT) *ScalarDeclNode {
	return &ScalarDeclNode{name, value, innerT}
}

func (node *ScalarDeclNode) Eval(scope *Scope) (Variable, error) {
	v, err := node.value.Eval(scope)
	if err != nil {
		return nil, err
	}
	if node.innerT != v.Type() {
		return nil, fmt.Errorf("%w: want=%s, got=%s", ErrScalarDeclTypesDiffer, TypeName(node.innerT), TypeName(v.Type()))
	}

	err = scope.ConstructScalar(node.name, v)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

type ArrayDeclNode struct {
	innerT VarT
	name   string
	sizes  []AstNode
	value  AstNode
}

func NewArrayDeclNode(innerType VarT, name string, sizes []AstNode, value AstNode) *ArrayDeclNode {
	return &ArrayDeclNode{innerType, name, sizes, value}
}

func (node *ArrayDeclNode) Eval(scope *Scope) (Variable, error) {
	v, err := node.value.Eval(scope)
	if err != nil {
		return nil, err
	}
	if v.Type() != node.innerT {
		return nil, fmt.Errorf("%w: want=%s, got=%s", ErrArrayAssignTypesDiffer, TypeName(node.innerT), TypeName(v.Type()))
	}
	sizes := make([]int, len(node.sizes))
	for i, sizeNode := range node.sizes {
		s, err := sizeNode.Eval(scope)
		if err != nil {
			return nil, err
		}
		if s.Type() != Int {
			return nil, ErrSizeNotInt
		}
		sizes[i] = s.(*Integer).Data()
	}
	err = scope.ConstructArray(node.name, sizes, v)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
