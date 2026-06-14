package interpreter

import (
	"errors"
	"fmt"
)

var ErrOperandNotScalar = errors.New("left operand is not scalar")

type AssignNode struct {
	left  AstNode
	right AstNode
}

func NewAssignNode(left, right AstNode) *AssignNode {
	return &AssignNode{left, right}
}

func (node *AssignNode) Eval(scope *Scope) (Variable, error) {
	l, err := node.left.Eval(scope)
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
