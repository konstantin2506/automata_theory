package interpreter

type MoveNode struct {
	count AstNode
}

func NewMoveNode(count AstNode) *MoveNode {
	return &MoveNode{count}
}

func (node *MoveNode) Eval(scope *Scope) (Variable, error) {
	v, err := node.count.Eval(scope)
	if err != nil {
		return nil, err
	}
	if v.Type() != Int {
		return nil, ErrVarInvalidType
	}
	n := v.(*Integer).Data()
	err = scope.globalScope.game.Move(n)
	return nil, err
}

type RotateNode struct {
	count AstNode
}

func NewRotateNode(count AstNode) *RotateNode {
	return &RotateNode{count}
}

func (node *RotateNode) Eval(scope *Scope) (Variable, error) {
	v, err := node.count.Eval(scope)
	if err != nil {
		return nil, err
	}
	if v.Type() != Int {
		return nil, ErrVarInvalidType
	}
	n := v.(*Integer).Data()
	scope.globalScope.game.Rotate(n)
	return nil, nil
}

type SurrNode struct{}

func NewSurrNode() *SurrNode {
	return &SurrNode{}
}

func (node *SurrNode) Eval(scope *Scope) (Variable, error) {
	sizes := []int{3, 5, 2}

	arr, err := NewArray(sizes, NewVariableBool(false))
	if err != nil {
		return nil, err
	}

	surr := scope.globalScope.game.Surroundings()
	for i := range 3 {
		for j := range 5 {
			for k := range 2 {
				err := arr.Assign([]int{i + 1, j + 1, k + 1}, NewVariableBool(surr[i][j][k]))
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return arr, nil
}
