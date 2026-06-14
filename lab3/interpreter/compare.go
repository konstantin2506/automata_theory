package interpreter

import "fmt"

type compT int

const (
	eq compT = iota
	lt
	lte
	gt
	gte
	eqMap
	ltMap
	lteMap
	gtMap
	gteMap
)

var compareNames = map[compT]string{
	eq:     "eq",
	lt:     "lt",
	lte:    "lte",
	gt:     "gt",
	gte:    "gte",
	eqMap:  "eq[]",
	ltMap:  "lt[]",
	lteMap: "lte[]",
	gtMap:  "gt[]",
	gteMap: "gte[]",
}

func CompareName(x compT) string {
	return compareNames[x]
}

type CompareNode struct {
	left      AstNode
	predicate func(int, int) bool
	compType  compT
	mapped    bool
}

func NewEqReduceNode(left, right AstNode) *CompareNode {
	return &CompareNode{left, func(l, r int) bool { return l == r }, eq, false}
}

func NewLtReduceNode(left, right AstNode) *CompareNode {
	return &CompareNode{left, func(l, r int) bool { return l < r }, lt, false}
}

func NewLteReduceNode(left, right AstNode) *CompareNode {
	return &CompareNode{left, func(l, r int) bool { return l <= r }, lte, false}
}

func NewGtReduceNode(left, right AstNode) *CompareNode {
	return &CompareNode{left, func(l, r int) bool { return l > r }, gt, false}
}

func NewGteReduceNode(left, right AstNode) *CompareNode {
	return &CompareNode{left, func(l, r int) bool { return l >= r }, gte, false}
}

func NewEqMapNode(left, right AstNode) *CompareNode {
	return &CompareNode{left, func(l, r int) bool { return l == r }, eq, true}
}

func NewLtMapNode(left, right AstNode) *CompareNode {
	return &CompareNode{left, func(l, r int) bool { return l < r }, lt, true}
}

func NewLteMapNode(left, right AstNode) *CompareNode {
	return &CompareNode{left, func(l, r int) bool { return l <= r }, lte, true}
}

func NewGtMapNode(left, right AstNode) *CompareNode {
	return &CompareNode{left, func(l, r int) bool { return l > r }, gt, true}
}

func NewGteMapNode(left, right AstNode) *CompareNode {
	return &CompareNode{left, func(l, r int) bool { return l >= r }, gte, true}
}

func CompReduce(v *VarArray, predicate func(int, int) bool) (Variable, error) {
	trueCount := 0
	for i := range v.data {
		l := (v.data[i]).(*Integer)
		resBool := predicate(l.Data(), 0)
		if resBool {
			trueCount++
		}
	}
	if trueCount > len(v.data)/2 {
		return NewVariableBool(true), nil
	}

	return NewVariableBool(false), nil
}

func CompMap(v *VarArray, predicate func(int, int) bool) (Variable, error) {
	res, err := NewArray(v.sizes, NewVariableBool(false))
	if err != nil {
		return nil, err
	}

	for i := range v.data {
		l := (v.data[i]).(*Integer)
		resBool := predicate(l.Data(), 0)
		res.data[i] = NewVariableBool(resBool)
	}

	return res, nil
}

func (node *CompareNode) Eval(scope *Scope) (Variable, error) {
	l, err := node.left.Eval(scope)
	if err != nil {
		return nil, err
	}

	lArr, okLarr := l.(*VarArray)
	if okLarr && lArr.InnerType() == Int {
		switch node.mapped {
		case true:
			return CompMap(lArr, node.predicate)
		case false:
			return CompReduce(lArr, node.predicate)
		}
	}

	lInt, okLeft := l.(*Integer)
	if !okLeft {
		return nil, fmt.Errorf("left %w in %s", ErrOperandIsNotBoolean, CompareName(node.compType))
	}
	result := node.predicate(lInt.Data(), 0)
	return NewVariableBool(result), nil
}
