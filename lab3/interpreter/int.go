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

func NewVariableInt(value int) Variable {
	return &Integer{value}
}
