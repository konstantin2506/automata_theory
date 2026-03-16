// Package nodes
// Ast nodes interface and impls
package nodes

type Node interface {
	Children() []Node
	Name() string
}
