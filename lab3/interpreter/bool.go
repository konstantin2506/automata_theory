package interpreter

import "fmt"

type Boolean struct {
	data bool
}

func (v *Boolean) Assign(right Variable) {
	rb := right.(*Boolean)
	v.data = rb.data
}

func (v *Boolean) Print() {
	fmt.Printf("%t", v.data)
}

func (v *Boolean) Type() VarT {
	return Bool
}

func (v *Boolean) Copy() Variable {
	return NewVariableBool(v.data)
}

func (v *Boolean) Data() bool {
	return v.data
}

func NewVariableBool(value bool) Variable {
	return &Boolean{value}
}

type BooleanNode struct {
	x bool
}

func NewBooleanNode(x bool) *BooleanNode {
	return &BooleanNode{x}
}

func (node *BooleanNode) Eval(scope *Scope) (Variable, error) {
	return NewVariableBool(node.x), nil
}
