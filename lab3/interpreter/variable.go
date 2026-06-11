package interpreter

type VarT int

const (
	Int VarT = iota
	Bool
)

type Variable interface {
	Type() VarT
}
