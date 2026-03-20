// Package ast
// Kostya's regular expressions library
package ast

type Ast struct {
	// root Node
}

//"abc(lq+)*"
//[a b c l]
//+ -> [+node]
//q ->

func BuildAst(str string) {
	bs := ByteStack{}
	ns := NodeStack{}

	for i := 0; i < len(str); i++ {
		ch := str[i]
		switch ch {
		case ')':
			for bs.Top() != '(' {
				ns = HandleByte(bs.Top(), ns)
				bs.Pop()
			}
			bs.Pop()
		case '(':
			bs.Push(ch)
		default:
			bs.Push(ch)
			ns.Push(&CharNode{ch})
		}
	}
	if !ns.isEmpty() {
		return
	}
}

func HandleByte(ch byte, ns NodeStack) NodeStack {
	switch ch {
	default:
		if !ns.isEmpty() {
			switch ns.Top().Type() {
			case Concat:
				node := ns.Top()
				ns.Pop()
				newNode := CharNode{ch}
				children := append(node.Children(), &newNode)
				ns.Push(&ConcatNode{children})
			case Char:
				node := ns.Top()
				ns.Pop()
				children := []Node{}
				children = append(children, node)
				children = append(children, &CharNode{ch})
				ns.Push(&ConcatNode{children})

			}
		} else {
			ns.Push(&CharNode{ch})
		}
	}
	return ns
}
