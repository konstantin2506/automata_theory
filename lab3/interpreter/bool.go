package interpreter

type Boolean struct {
	data bool
}

func (v *Boolean) Assign(value bool) {
	v.data = value
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
