package interpreter

import "fmt"

type MostNode struct {
	child AstNode
}

func NewMostNode(child AstNode) *MostNode {
	return &MostNode{child}
}

func (node *MostNode) Eval(scope *Scope) (Variable, error) {
	v, err := node.child.Eval(scope)
	if err != nil {
		return nil, err
	}
	if v.Type() == Array && v.(*VarArray).InnerType() == Bool {
		cnt := 0
		for _, elem := range v.(*VarArray).data {
			if elem.(*Boolean).Data() {
				cnt++
			}
		}
		if cnt > len(v.(*VarArray).data)/2 {
			return NewVariableBool(true), nil
		}
		return NewVariableBool(false), nil
	}
	if v.Type() == Bool {
		return NewVariableBool(v.(*Boolean).Data()), nil
	}
	return nil, fmt.Errorf("%w in most", ErrVarInvalidType)
}
