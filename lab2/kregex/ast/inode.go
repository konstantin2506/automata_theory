package ast

type Node interface {
	Children() []Node
	Type() NodeT
	String() string
}

type NodeT int

const (
	Concat NodeT = iota
	Kleene
	NamedGroup
	Optional
	Or
	Reference
	Repeat
	Char
	OpenParen
)
