package ast

type Regex struct {
	nfa      Nfa
	dfa      DFA
	tree     Ast
	compiled bool
}

// builds only nfa for expr
func NewRegex(expr string) (Regex, error) {
	tree, err := BuildAst(expr)
	if err != nil {
		return Regex{}, err
	}
	nfa := BuildFromAst(tree.GetRoot())
	return Regex{
		nfa:      nfa,
		tree:     tree,
		compiled: false,
	}, nil
}

func Compile(expr string) (Regex, error) {
	regex, err := NewRegex(expr)
	if err != nil {
		return Regex{}, err
	}
	specMap := NewNodeSpecMap()
	tree := regex.tree
	chars, charNums := MarkChars(&tree)

	ComputeNullable(&tree, specMap, charNums)
	first := ComputeFirst(&tree, specMap, charNums)
	last := ComputeLast(&tree, specMap, charNums)

	follow := ComputeFollow(specMap)
	dfa := BuildDFA(tree.GetRoot(), follow, first, last, chars, charNums, specMap)
	regex.compiled = true
	regex.dfa = dfa
	return regex, nil
}

// if compiled uses dfa, else nfa
func Search(str string, regex *Regex) bool {
	found := false
	switch regex.compiled {
	case true:
		//_, found = regex.nfa.Search(str)
	case false:
		//found = regex.nfa.Search(str)
	}
	return found
}

//func SearchWithGroups(str string, regex *Regex) (map[int]string, bool)
