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

func NewVariableBool(value bool) Variable {
	return &Boolean{value}
}
