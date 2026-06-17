package interpreter

import "errors"

var (
	ErrCastVecToScalar   = errors.New("trying to cast vector type to scalar type")
	ErrCastToUnknownType = errors.New("trying to cast to unknown type")
)

type CastNode struct {
	child  AstNode
	toType VarT
}

func NewToIntegerCastNode(child AstNode) *CastNode {
	return &CastNode{child, Int}
}

func NewToBooleanCastNode(child AstNode) *CastNode {
	return &CastNode{child, Bool}
}

func convBoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func convIntToBool(i int) bool {
	return i != 0
}

func (node *CastNode) Eval(scope *Scope) (Variable, error) {
	from, err := node.child.Eval(scope)
	if err != nil {
		return nil, err
	}
	if from.Type() == Array {
		return nil, ErrCastVecToScalar
	}

	switch node.toType {
	case Int:
		res := 0
		switch typed := from.(type) {
		case *Integer:
			res = typed.Data()
		case *Boolean:
			res = convBoolToInt(typed.Data())
		}
		return NewVariableInt(res), nil
	case Bool:
		res := false
		switch typed := from.(type) {
		case *Integer:
			res = convIntToBool(typed.Data())
		case *Boolean:
			res = typed.Data()
		}
		return NewVariableBool(res), nil
	}
	return nil, ErrCastToUnknownType
}
