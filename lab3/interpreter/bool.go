package interpreter

type Boolean struct {
	data bool
}

func (v *Boolean) Assign(right Variable) {
	rb := right.(*Boolean)
	v.data = rb.data
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
