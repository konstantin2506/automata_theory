package interpreter

type VarT int

const (
	Int VarT = iota
	Bool
	Array
)

var typeNames = map[VarT]string{
	Int:   "Integer",
	Bool:  "Boolean",
	Array: "Array",
}

func TypeName(t VarT) string {
	return typeNames[t]
}

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
