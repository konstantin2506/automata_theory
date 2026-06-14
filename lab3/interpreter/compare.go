package interpreter

import "fmt"

type compT int

const (
	eq compT = iota
	lt
	lte
	gt
	gte
)

var compareNames = map[compT]string{
	eq:  "==",
	lt:  "<",
	lte: "<=",
	gt:  ">",
	gte: ">=",
}

func CompareName(x compT) string {
	return compareNames[x]
}

type CompareNode struct {
	left      AstNode
	right     AstNode
	predicate func(int, int) bool
	compType  compT
}

func NewEqNode(left, right AstNode) AstNode {
	return &CompareNode{left, right, func(l, r int) bool { return l == r }, eq}
}

func NewLtNode(left, right AstNode) AstNode {
	return &CompareNode{left, right, func(l, r int) bool { return l < r }, lt}
}

func NewLteNode(left, right AstNode) AstNode {
	return &CompareNode{left, right, func(l, r int) bool { return l <= r }, lte}
}

func NewGtNode(left, right AstNode) AstNode {
	return &CompareNode{left, right, func(l, r int) bool { return l > r }, gt}
}

func NewGteNode(left, right AstNode) AstNode {
	return &CompareNode{left, right, func(l, r int) bool { return l >= r }, gte}
}

func CompReduce(v, other *VarArray, predicate func(int, int) bool) (Variable, error) {
	err := v.CmpTypeWith(other)
	if err != nil {
		return nil, err
	}
	trueCount := 0
	for i := range v.data {
		l := (v.data[i]).(*Integer)
		r := (other.data[i]).(*Integer)
		resBool := predicate(l.Data(), r.Data())
		if resBool {
			trueCount++
		}
	}
	if trueCount > len(v.data)/2 {
		return NewVariableBool(true), nil
	}

	return NewVariableBool(false), nil
}

func CompMap(v, other *VarArray, predicate func(int, int) bool) (Variable, error) {
	err := v.CmpTypeWith(other)
	if err != nil {
		return nil, err
	}
	res, err := NewArray(v.sizes, NewVariableBool(false))
	if err != nil {
		return nil, err
	}

	for i := range v.data {
		l := (v.data[i]).(*Integer)
		r := (other.data[i]).(*Integer)
		resBool := predicate(l.Data(), r.Data())
		res.(*VarArray).data[i] = NewVariableBool(resBool)
	}

	return res, nil
}

func (node *CompareNode) Eval(scope *Scope) (Variable, error) {
	l, err := node.left.Eval(scope)
	if err != nil {
		return nil, err
	}
	r, err := node.right.Eval(scope)
	if err != nil {
		return nil, err
	}
	lInt, okLeft := l.(*Integer)
	if !okLeft {
		return nil, fmt.Errorf("left %w in %s", ErrOperandIsNotBoolean, CompareName(node.compType))
	}
	rInt, okRight := r.(*Integer)
	if !okRight {
		return nil, fmt.Errorf("right %w in %s", ErrOperandIsNotBoolean, CompareName(node.compType))
	}
	result := node.predicate(lInt.Data(), rInt.Data())
	return NewVariableBool(result), nil
}
