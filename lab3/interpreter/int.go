package interpreter

type Integer struct {
	data int
}

func (v *Integer) Assign(value int) {
	v.data = value
}

func (v *Integer) Type() VarT {
	return Int
}

func (v *Integer) Copy() Variable {
	return NewVariableInt(v.data)
}

func NewVariableInt(value int) Variable {
	return &Integer{value}
}

func (v *Integer) Data() int {
	return v.data
}
