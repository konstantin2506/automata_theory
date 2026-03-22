// Package ast
// Kostya's regular expressions library
package ast

import (
	"fmt"
	"strings"
)

type Ast struct {
	root Node
}

//"abc(lq+)*"
//[a b c l]
//+ -> [+node]
//q ->

func BuildAst(str string) (Ast, error) {
	ns := NodeStack{}

	for i := 0; i < len(str); i++ {
		ch := str[i]
		switch ch {
		case ')':
			ns = HandleNodes(ns)
		case '(':
			ns.Push(&OpenParenNode{})
		case '?':
			ns.Push(&OptionalNode{})
		default:
			ns.Push(&CharNode{ch})
		}
	}
	if ns.isEmpty() {
		return Ast{}, fmt.Errorf("error in stack while BuildAst")
	}
	return Ast{ns.Top()}, nil
}

func HandleNodes(ns NodeStack) NodeStack {
	newNodeStack := NodeStack{}

	for ns.Top().Type() != OpenParen {
		current := ns.Top()
		ns.Pop()

		switch current.Type() {
		case Optional:
			if ns.Top().Type() == OpenParen {
				newNodeStack.Push(current)
				break
			}
			child := ns.Top()
			ns.Pop()
			current.(*OptionalNode).child = child
			newNodeStack.Push(current)
		case Char:
			newNodeStack.Push(current)
		case Concat:
			newNodeStack.Push(current)
		}
	}
	ns.Pop()
	for newNodeStack.Size() != 1 {
		current := newNodeStack.Top()
		newNodeStack.Pop()
		right := newNodeStack.Top()
		newNodeStack.Pop()
		newChildren := []Node{}
		if current.Type() == Concat {
			newChildren = append(current.(*ConcatNode).children, right)
		} else {
			newChildren = append(newChildren, current)
			newChildren = append(newChildren, right)
		}
		newNodeStack.Push(&ConcatNode{newChildren})
	}

	ns.Push(newNodeStack.Top())

	return ns
}

func (ast *Ast) Print(depth int) {
	if depth == 1 {
		fmt.Printf("%s\n", ast.root.String())
	}
	for _, child := range ast.root.Children() {
		if child != nil {
			fmt.Printf("%s", strings.Repeat(" ", depth*2))
			fmt.Printf("%s\n", child.String())
			ast := Ast{child}
			ast.Print(depth + 1)
		} else {
			fmt.Printf("nil\n")
		}
	}
}
