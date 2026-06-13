package interpreter

type StatementListNode struct {
	statements []AstNode
}

func NewStatementListNode() AstNode {
	return &StatementListNode{[]AstNode{}}
}

func (node *StatementListNode) PushStatement(statement AstNode) {
	node.statements = append(node.statements, statement)
}

func (node *StatementListNode) Eval(scope *Scope) (Variable, error) {
	for _, statement := range node.statements {
		res, err := statement.Eval(scope)
		if err != nil {
			return nil, err
		}
		if res != nil {
			return res, nil
		}
	}
	return nil, nil
}
