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

func (ast *Ast) GetRoot() Node {
	return ast.root
}

func BuildAst(str string) (Ast, error) {
	str = fmt.Sprintf("(%s)", str)
	ns := NodeStack{}

	for i := 0; i < len(str); i++ {
		ch := str[i]
		switch ch {
		case ')':
			ns = HandleParens(ns)
		case '(':
			if i+1 < len(str) && str[i+1] == ')' {
				return Ast{}, fmt.Errorf("empty group not allowed at position %d", i)
			}
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
				ns.Push(&CharNode{ch, 0})
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
			ns.Push(&CharNode{str[i+1], 0})
			i += 2
		case '<':
			if i == 0 || str[i-1] != '(' {
				return Ast{}, fmt.Errorf("group name does not start with '(', pos=%d", i)
			}
			start := i
			end := strings.Index(str[i:], ">")
			if end != -1 {
				if len(str) > start+end && str[start+end+1] == ')' {
					return Ast{}, fmt.Errorf("empty named group, pos: %d", start)
				}
				ns.Push(&NamedGroupNode{nil, str[start+1 : start+end]})
				i += end
			} else {
				return Ast{}, fmt.Errorf("'<' does not have a paired '>'")
			}
		case '>':
			return Ast{}, fmt.Errorf("'>' does not have a paired '<'")

		default:
			ns.Push(&CharNode{ch, 0})
		}
	}

	if ns.Size() != 1 {
		return Ast{}, fmt.Errorf("error in stack while BuildAst (bad parens)")
	}
	return Ast{ns.Top()}, nil
}

func HandleParens(ns NodeStack) NodeStack {
	concatStack, gname := HandleFirstPriorityOps(&ns)

	orNode := ConcatinateNodes(&concatStack)
	HandleLastNode(&ns, &concatStack, gname, orNode)
	return ns
}

func (ast *Ast) Print(depth int) {
	fmt.Println(ast.TraverseRLR("", 1, '\n', " "))
}

func (ast *Ast) TraverseRLRSpace() string {
	return ast.TraverseRLR("", 1, ' ', "")
}

func (ast *Ast) TraverseRLR(str string, depth int, delim byte, depthString string) string {
	if depth == 1 {
		str = fmt.Sprintf("%s%c", ast.root.String(), delim)
	}
	for _, child := range ast.root.Children() {
		if child != nil {
			str += strings.Repeat(depthString, depth*2)
			str += fmt.Sprintf("%s%c", child.String(), delim)
			ast := Ast{child}
			str = ast.TraverseRLR(str, depth+1, delim, depthString)
		} else {
			str += fmt.Sprintf("nil%c", delim)
		}
	}
	return str
}

func NewSubAst(root Node) Ast {
	return Ast{root}
}

func HandleFirstPriorityOps(ns *NodeStack) (NodeStack, Node) {
	concatStack := NodeStack{}
	var gname Node
	for ns.Top().Type() != OpenParen { // связываем операции первого приоритета, добавляем в стек для конкатенации
		current := ns.Top()
		ns.Pop()

		switch current.Type() {
		case NamedGroup:
			if current.(*NamedGroupNode).child == nil {
				gname = current
			} else {
				concatStack.Push(current)
			}
		case Or:
			concatStack.Push(current)
		case Optional:
			if ns.Top().Type() == OpenParen {
				concatStack.Push(current)
				break
			}
			child := ns.Top()
			ns.Pop()
			current.(*OptionalNode).child = child
			concatStack.Push(current)
		case Char:
			concatStack.Push(current)
		case Concat:
			concatStack.Push(current)
		case Kleene:
			if ns.Top().Type() == OpenParen {
				concatStack.Push(current)
				break
			}
			if current.(*KleeneNode).child == nil {
				child := ns.Top()
				ns.Pop()
				current.(*KleeneNode).child = child
			}
			concatStack.Push(current)
		case Repeat:
			if ns.Top().Type() == OpenParen {
				concatStack.Push(current)
				break
			}
			child := ns.Top()
			ns.Pop()
			current.(*RepeatNode).child = child
			concatStack.Push(current)
		}
	}
	ns.Pop()
	return concatStack, gname
}

func ConcatinateNodes(concatStack *NodeStack) OrNode {
	orNode := OrNode{}
	for concatStack.Size() != 1 {
		current := concatStack.Top()
		concatStack.Pop()
		right := concatStack.Top()
		concatStack.Pop()
		newChildren := []Node{}
		if right.Type() == Or && !concatStack.isEmpty() {
			orNode.childs = append(orNode.childs, current)
			if len(right.Children()) != 0 {
				orNode.childs = append(orNode.childs, right)
			}
		} else {
			if current.Type() == Concat {
				newChildren = append(current.(*ConcatNode).children, right)
			} else {
				newChildren = append(newChildren, current)
				newChildren = append(newChildren, right)
			}
			concatStack.Push(&ConcatNode{newChildren})
		}
	}
	return orNode
}

func HandleLastNode(ns, concatStack *NodeStack, gname Node, orNode OrNode) {
	if gname == nil {
		if len(orNode.childs) == 0 {
			ns.Push(concatStack.Top())
		} else {
			orNode.childs = append(orNode.childs, concatStack.Top())
			ns.Push(&orNode)
		}
	} else {
		if len(orNode.childs) == 0 {
			gname.(*NamedGroupNode).child = concatStack.Top()
		} else {
			orNode.childs = append(orNode.childs, concatStack.Top())
			gname.(*NamedGroupNode).child = &orNode
		}
		ns.Push(gname)
	}
}
