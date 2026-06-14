package interpreter

type Integer struct {
	data int
}

func (v *Integer) Assign(value Variable) {
	r := value.(*Integer)
	v.data = r.data
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

type IntegerNode struct {
	x int
}

func NewIntegerNode(x int) AstNode {
	return &IntegerNode{x}
}

func (node *IntegerNode) Eval(scope *Scope) (Variable, error) {
	return NewVariableInt(node.x), nil
}
