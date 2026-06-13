package interpreter

import (
	"errors"
	"fmt"
)

type BinBoolOpT int

const (
	And BinBoolOpT = iota
	Or
)

var ErrOperandIsNotBoolean = errors.New("operand must be boolean")

var boolBinOpNames = map[BinBoolOpT]string{
	And: "and",
	Or:  "or",
}

func BinOpName(op BinBoolOpT) string {
	return boolBinOpNames[op]
}

type BinBoolOpNode struct {
	left      AstNode
	right     AstNode
	predicate func(l, r bool) bool
	opType    BinBoolOpT
}

func NewAndNode(left, right AstNode) AstNode {
	return &BinBoolOpNode{left, right, func(l, r bool) bool { return l && r }, And}
}

func NewOrNode(left, right AstNode) AstNode {
	return &BinBoolOpNode{left, right, func(l, r bool) bool { return l || r }, Or}
}

func (node *BinBoolOpNode) Eval(scope *Scope) (Variable, error) {
	l, err := node.left.Eval(scope)
	if err != nil {
		return nil, err
	}
	r, err := node.right.Eval(scope)
	if err != nil {
		return nil, err
	}
	lBool, okLeft := l.(*Boolean)
	if !okLeft {
		return nil, fmt.Errorf("left %w in %s", ErrOperandIsNotBoolean, BinOpName(node.opType))
	}
	rBool, okRight := r.(*Boolean)
	if !okRight {
		return nil, fmt.Errorf("right %w in %s", ErrOperandIsNotBoolean, BinOpName(node.opType))
	}
	result := node.predicate(lBool.Data(), rBool.Data())
	return NewVariableBool(result), nil
}
