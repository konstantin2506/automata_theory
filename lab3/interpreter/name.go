package interpreter

type NameNode struct {
	name string
}

func NewNameNode(name string) *NameNode {
	return &NameNode{name}
}

func (node *NameNode) Eval(scope *Scope) (Variable, error) {
	return scope.FindVariableDepth(node.name)
}
