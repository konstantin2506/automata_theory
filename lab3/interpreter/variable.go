package interpreter

type VarT int

const (
	Int VarT = iota
	Bool
	Array
)

type Variable interface {
	Type() VarT
	Copy() Variable
}

type Scalar interface {
	Assign(value Variable)
}

type Vector interface {
	Assign(indices []int, value Variable) error
	InnerType() VarT
}
