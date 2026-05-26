package kregex

type Regex struct {
	nfa      Nfa
	dfa      Dfa
	tree     Ast
	compiled bool
}

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
	dfa := Minimize(NewDfa((regex.tree)))

	regex.compiled = true
	regex.dfa = dfa
	return regex, nil
}

func Search(str string, regex *Regex) string {
	var found string
	switch regex.compiled {

	case true:
		found = regex.dfa.Search(str)
	case false:
		found, _ = regex.nfa.Search(str)
	}
	return found
}

func SearchWithGroups(str string, regex *Regex) (string, map[string]string) {
	return regex.nfa.Search(str)
}
