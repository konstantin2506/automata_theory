package interpreter

type SizeNode struct {
	ofWhat AstNode
}

func NewSizeNode(ofWhat AstNode) *SizeNode {
	return &SizeNode{ofWhat}
}

func (node *SizeNode) Eval(scope *Scope) (Variable, error) {
	v, err := node.ofWhat.Eval(scope)
	if err != nil {
		return nil, err
	}
	if v.Type() != Array {
		res, err := NewArray([]int{1}, NewVariableInt(1))
		return res, err
	}
	sizes := v.(*VarArray).sizes
	res, err := NewArray([]int{len(sizes)}, NewVariableInt(0))
	if err != nil {
		return nil, err
	}
	for i := range sizes {
		err := res.Assign([]int{i + 1}, NewVariableInt(sizes[i]))
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
