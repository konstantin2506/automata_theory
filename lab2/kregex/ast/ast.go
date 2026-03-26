// Package ast
// Kostya's regular expressions library
package ast

import (
	"fmt"
	"strconv"
	"strings"
)

type Ast struct {
	root Node
}

// "(ab(c?)...)"
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
			if i+3 < len(str) && str[i+1:i+4] == "..." {
				return Ast{}, fmt.Errorf("combo '?...' is incorrect")
			}
			ns.Push(&OptionalNode{})
		case '.':
			if (i < len(str)-2) && str[i:i+3] == "..." {
				if i+3 < len(str) && str[i+3] == '?' {
					return Ast{}, fmt.Errorf("combo '...?' is incorrect")
				}
				ns.Push(&KleeneNode{})
				i += 2
			} else {
				ns.Push(&CharNode{ch})
			}
		case '{':
			start := i
			end := strings.Index(str[i:], "}")
			if end != -1 {
				count, err := strconv.ParseUint(str[start+1:start+end], 10, 64)
				if err != nil {
					return Ast{}, fmt.Errorf("incorrect count in braces: %s", str[start+1:start+end])
				}
				ns.Push(&RepeatNode{nil, uint32(count)})
				i += end
			} else {
				return Ast{}, fmt.Errorf("'{' does not have a paired '}'")
			}
		case '}':
			return Ast{}, fmt.Errorf("'}' does not have a paired '{'")
		case '|':
			ns.Push(&OrNode{})
		case '%':
			if i > len(str)-2 || str[i+2] != '%' {
				return Ast{}, fmt.Errorf("incorrect escape symbol usage, pos: %d", i)
			}
			ns.Push(&CharNode{str[i+1]})
			i += 2
		case '<':
			if i == 0 || str[i-1] != '(' {
				return Ast{}, fmt.Errorf("group name does not start with '(', pos=%d", i)
			}
			start := i
			end := strings.Index(str[i:], ">")
			if end != -1 {
				ns.Push(&NamedGroupNode{nil, str[start+1 : start+end]})
				i += end
			} else {
				return Ast{}, fmt.Errorf("'<' does not have a paired '>'")
			}
		case '>':
			return Ast{}, fmt.Errorf("'>' does not have a paired '<'")

		default:
			ns.Push(&CharNode{ch})
		}
	}

	if ns.Size() != 1 {
		return Ast{}, fmt.Errorf("error in stack while BuildAst (bad parens)")
	}
	return Ast{ns.Top()}, nil
}

func HandleNodes(ns NodeStack) NodeStack {
	newNodeStack := NodeStack{}
	withNamedGroup := false
	var gname Node
	for ns.Top().Type() != OpenParen {
		current := ns.Top()
		ns.Pop()

		switch current.Type() {
		case NamedGroup:

			if current.(*NamedGroupNode).child == nil {
				withNamedGroup = true
				gname = current
			} else {
				newNodeStack.Push(current)
			}
		case Or:
			newNodeStack.Push(current)
		case Optional:
			if ns.Top().Type() == OpenParen {
				newNodeStack.Push(current)
				break
			}
			if ns.Top().Type() == NamedGroup {
				newNodeStack.Push(current)
				withNamedGroup = true
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
		case Kleene:
			if ns.Top().Type() == OpenParen {
				newNodeStack.Push(current)
				break
			}
			if ns.Top().Type() == NamedGroup {
				newNodeStack.Push(current)
				withNamedGroup = true
				break
			}
			child := ns.Top()
			ns.Pop()
			current.(*KleeneNode).child = child
			newNodeStack.Push(current)

		case Repeat:
			if ns.Top().Type() == OpenParen {
				newNodeStack.Push(current)
				break
			}
			if ns.Top().Type() == NamedGroup {
				newNodeStack.Push(current)
				withNamedGroup = true
				break
			}
			child := ns.Top()
			ns.Pop()
			current.(*RepeatNode).child = child
			newNodeStack.Push(current)
		}
	}
	ns.Pop()
	if newNodeStack.isEmpty() {
		if withNamedGroup {
			gname.(*NamedGroupNode).child = ns.Top()
			ns.Pop()
		}
		ns.Push(gname)
		return ns
	}
	for newNodeStack.Size() != 1 {
		current := newNodeStack.Top()
		newNodeStack.Pop()
		right := newNodeStack.Top()
		newNodeStack.Pop()
		newChildren := []Node{}
		if right.Type() == Or {
			if current.Type() == Or {
				current.(*OrNode).childs = append(current.(*OrNode).childs, newNodeStack.Top())
				newNodeStack.Pop()
				newNodeStack.Push(current)
			} else {
				childs := []Node{}
				childs = append(childs, current)
				childs = append(childs, newNodeStack.Top())
				newNodeStack.Pop()
				right.(*OrNode).childs = childs
				newNodeStack.Push(right)
			}
		} else {
			if current.Type() == Concat {
				newChildren = append(current.(*ConcatNode).children, right)
			} else {
				newChildren = append(newChildren, current)
				newChildren = append(newChildren, right)
			}
			newNodeStack.Push(&ConcatNode{newChildren})
		}

	} // TODO
	if withNamedGroup && gname != nil {
		gname.(*NamedGroupNode).child = newNodeStack.Top()
		ns.Push(gname)
	} else if withNamedGroup && gname == nil {
		ns.Push(newNodeStack.Top())
	}
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
