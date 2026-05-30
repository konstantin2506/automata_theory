package main

const (
	TypeDigit
)

type Value struct {
	Type
}

type Node interface {
	Eval() *int
}

type BinOpNode struct {
	left  Node
	right Node
}

func (n *BinOpNode) Eval() *int
