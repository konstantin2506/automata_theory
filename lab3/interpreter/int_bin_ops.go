package interpreter

import (
	"errors"
	"fmt"
)

type BinIntOpT int

const (
	Div BinIntOpT = iota
	Mul
	Sub
	Add
)

var (
	ErrOperandIsNotInteger = errors.New("operand must be integer")
	ErrZeroDivision        = errors.New("right operand is equal to zero")
)

var intBinOpNames = map[BinIntOpT]string{
	Div: "div",
	Mul: "mul",
	Sub: "sub",
	Add: "add",
}

func BinIntOpName(op BinIntOpT) string {
	return intBinOpNames[op]
}

type BinIntOpNode struct {
	left      AstNode
	right     AstNode
	predicate func(l, r int) int
	opType    BinIntOpT
}

func NewDivNode(left, right AstNode) AstNode {
	return &BinIntOpNode{left, right, func(l, r int) int { return l / r }, Div}
}

func NewMulNode(left, right AstNode) AstNode {
	return &BinIntOpNode{left, right, func(l, r int) int { return l * r }, Mul}
}

func NewSubNode(left, right AstNode) AstNode {
	return &BinIntOpNode{left, right, func(l, r int) int { return l - r }, Sub}
}

func NewAddNode(left, right AstNode) AstNode {
	return &BinIntOpNode{left, right, func(l, r int) int { return l + r }, Add}
}

func (node *BinIntOpNode) Eval(scope *Scope) (Variable, error) {
	l, err := node.left.Eval(scope)
	if err != nil {
		return nil, err
	}
	r, err := node.right.Eval(scope)
	if err != nil {
		return nil, err
	}
	lInt, okLeft := l.(*Integer)
	if !okLeft {
		return nil, fmt.Errorf("left %w in %s", ErrOperandIsNotInteger, BinIntOpName(node.opType))
	}
	rInt, okRight := r.(*Integer)
	if !okRight {
		return nil, fmt.Errorf("right %w in %s", ErrOperandIsNotInteger, BinIntOpName(node.opType))
	}
	if rInt.Data() == 0 && node.opType == Div {
		return nil, ErrZeroDivision
	}
	result := node.predicate(lInt.Data(), rInt.Data())
	return NewVariableInt(result), nil
}
