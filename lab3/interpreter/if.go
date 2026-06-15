package interpreter

import (
	"errors"
	"fmt"
)

var ErrConditionNotBool = errors.New("condition result is not bool")

type IfStatementNode struct {
	condition AstNode
	ifTrue    AstNode
	ifFalse   AstNode
}

func NewIfStatementNode(condition, ifTrue, ifFalse AstNode) *IfStatementNode {
	return &IfStatementNode{condition, ifTrue, ifFalse}
}

func (node *IfStatementNode) Eval(scope *Scope) (Variable, error) {
	res, err := node.condition.Eval(scope)
	if err != nil {
		return nil, err
	}
	if res.Type() != Bool {
		return nil, fmt.Errorf("%w: type=%s", ErrConditionNotBool, TypeName(res.Type()))
	}
	condition := res.(*Boolean).Data()
	childScope := NewScope(scope)
	if condition {
		resTrue, err := node.ifTrue.Eval(&childScope)
		if err != nil {
			return nil, err
		}
		return resTrue, nil
	} else {
		if node.ifFalse == nil {
			return nil, nil
		}
		resFalse, err := node.ifFalse.Eval(&childScope)
		if err != nil {
			return nil, err
		}
		return resFalse, nil
	}
}
