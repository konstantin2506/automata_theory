package interpreter

import (
	"errors"
	"fmt"
)

var ErrOperandNotScalar = errors.New("left operand is not scalar")

type AssignNode struct {
	name  string
	right AstNode
}

func NewAssignNode(name string, right AstNode) *AssignNode {
	return &AssignNode{name, right}
}

func (node *AssignNode) Eval(scope *Scope) (Variable, error) {
	l, err := scope.FindVariableDepth(node.name)
	if err != nil {
		return nil, err
	}
	r, err := node.right.Eval(scope)
	if err != nil {
		return nil, err
	}
	lScalar, ok := l.(Scalar)
	if !ok {
		return nil, fmt.Errorf("left %w", ErrOperandNotScalar)
	}
	_, ok = r.(Scalar)
	if !ok {
		return nil, fmt.Errorf("right %w", ErrOperandNotScalar)
	}
	lScalar.Assign(r)
	return nil, nil
}

type ArrayAssignNode struct {
	elem  *ArrayElemNode
	right AstNode
}

func NewArrayAssignNode(elem *ArrayElemNode, right AstNode) *ArrayAssignNode {
	return &ArrayAssignNode{elem, right}
}

func (node *ArrayAssignNode) Eval(scope *Scope) (Variable, error) {
	el, err := node.elem.Eval(scope)
	if err != nil {
		return nil, err
	}
	r, err := node.right.Eval(scope)
	if err != nil {
		return nil, err
	}
	if el.Type() != r.Type() {
		return nil, ErrArrayAssignTypesDiffer
	}
	lScalar, ok := el.(Scalar)
	if !ok {
		return nil, fmt.Errorf("left %w", ErrOperandNotScalar)
	}
	_, ok = r.(Scalar)
	if !ok {
		return nil, fmt.Errorf("right %w", ErrOperandNotScalar)
	}
	lScalar.Assign(r)
	return nil, nil
}
