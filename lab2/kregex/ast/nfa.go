package ast

import (
	"fmt"
	"os"
	"strings"
)

type NfaNodeT int

const (
	Default NfaNodeT = iota
	GroupStart
	GroupEnd
)

type EpsilonTransition struct {
	dst *NfaNode
}

type CharTransition struct {
	char byte
	dst  *NfaNode
}

type NfaNode struct {
	epsilonTransitions []EpsilonTransition
	charTransitions    []CharTransition
	nodeType           NfaNodeT
	groupName          string
}

func (n *NfaNode) addEpsilonWith(dst *NfaNode) {
	n.epsilonTransitions = append(n.epsilonTransitions, EpsilonTransition{dst})
}
func (n *NfaNode) addCharWith(src *NfaNode, char byte) {
	n.charTransitions = append(n.charTransitions, CharTransition{char, src})
}

type Nfa struct {
	start *NfaNode
	end   *NfaNode
}

func newNfaNode() *NfaNode {
	return &NfaNode{epsilonTransitions: make([]EpsilonTransition, 0),
		charTransitions: make([]CharTransition, 0),
		nodeType:        Default,
	}
}

func buildCharNfa(char byte) Nfa {
	first := newNfaNode()
	second := newNfaNode()

	transtion := CharTransition{char, second}
	first.charTransitions = append(first.charTransitions, transtion)
	return Nfa{first, second}
}

/*
  - ломает группы почему то хз
    func buildConcatNfa(nfas []Nfa) Nfa {
    if len(nfas) < 2 {
    panic("incorrect concat in buildConcatNfa")
    }
    result := nfas[0]
    for i := 1; i < len(nfas); i++ {
    next := nfas[i]
    //result.end.addEpsilonWith(next.start)
    result.end.epsilonTransitions = append(result.end.epsilonTransitions, next.start.epsilonTransitions...)
    result.end.charTransitions = append(result.end.charTransitions, next.start.charTransitions...)
    result.end = next.end
    }
    return result

}
*/
func buildConcatNfa(nfas []Nfa) Nfa {
	if len(nfas) < 2 {
		panic("incorrect concat in buildConcatNfa")
	}
	result := nfas[0]
	for i := 1; i < len(nfas); i++ {
		next := nfas[i]
		result.end.addEpsilonWith(next.start)
		result.end = next.end
	}
	return result
}

func buildKleeneNfa(current Nfa) Nfa {
	first := newNfaNode()
	last := newNfaNode()
	first.addEpsilonWith(current.start)
	first.addEpsilonWith(last)
	current.end.addEpsilonWith(last)
	current.end.addEpsilonWith(current.start)

	return Nfa{first, last}
}

func buildNamedGroup(current Nfa, name string) Nfa {
	first := newNfaNode()
	last := newNfaNode()
	first.nodeType = GroupStart
	first.groupName = name
	last.nodeType = GroupEnd
	last.groupName = name

	first.addEpsilonWith(current.start)
	current.end.addEpsilonWith(last)
	return Nfa{first, last}
}

func buildOrNfa(nfas []Nfa) Nfa {
	first := newNfaNode()
	last := newNfaNode()
	for _, nfa := range nfas {
		first.addEpsilonWith(nfa.start)
		nfa.end.addEpsilonWith(last)
	}
	return Nfa{first, last}
}

func buildOptionalNfa(nfa Nfa) Nfa {
	first := newNfaNode()
	last := newNfaNode()
	first.addEpsilonWith(nfa.start)
	first.addEpsilonWith(last)
	nfa.end.addEpsilonWith(last)
	return Nfa{first, last}
}

func BuildFromAst(root Node) Nfa {
	if root == nil {
		return Nfa{nil, nil}
	}
	switch root := root.(type) {
	case *ConcatNode:
		nfas := []Nfa{}
		for _, child := range root.Children() {
			nfas = append(nfas, BuildFromAst(child))
		}
		return buildConcatNfa(nfas)

	case *OrNode:
		nfas := []Nfa{}
		for _, child := range root.Children() {
			nfas = append(nfas, BuildFromAst(child))
		}
		return buildOrNfa(nfas)

	case *CharNode:
		return buildCharNfa(root.char)
	case *NamedGroupNode:
		return buildNamedGroup(BuildFromAst(root.child), root.name)
	case *KleeneNode:
		return buildKleeneNfa(BuildFromAst(root.Children()[0]))
	case *OptionalNode:
		return buildOptionalNfa(BuildFromAst(root.child))
	default:
		panic("unknown Node type")
	}

}

func (this *NfaNode) buildEpsilonClosure() map[*NfaNode]struct{} {
	closure := map[*NfaNode]struct{}{}
	stack := []*NfaNode{}
	stack = append(stack, this)

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		for _, transition := range current.epsilonTransitions {
			next := transition.dst
			_, ok := closure[next]
			if !ok {
				closure[next] = struct{}{}
				stack = append(stack, next)
			}
		}
	}
	return closure

}

func (nfa *Nfa) Search(str string) (bool, map[string]string) {
	if nfa.start == nil && nfa.end == nil && str == "" {
		return true, nil
	}
	currentStates := nfa.start.buildEpsilonClosure()
	currentStates[nfa.start] = struct{}{}
	groupsIndexs := map[string][]int{}
	groups := map[string]string{}

	for i := 0; i < len(str); i++ {
		ch := str[i]
		nextStates := map[*NfaNode]struct{}{}

		for state := range currentStates {
			if state.nodeType == GroupStart {
				if _, ok := groupsIndexs[state.groupName]; !ok {
					groupsIndexs[state.groupName] = make([]int, 2)
					groupsIndexs[state.groupName][0] = -1
					groupsIndexs[state.groupName][1] = -1
				}
				groupsIndexs[state.groupName][0] = i
			} else if state.nodeType == GroupEnd {
				if _, ok := groupsIndexs[state.groupName]; !ok {
					groupsIndexs[state.groupName] = make([]int, 2)
					groupsIndexs[state.groupName][0] = -1
					groupsIndexs[state.groupName][1] = -1
				}
				groupsIndexs[state.groupName][1] = i - 1

			}

			for _, transition := range state.charTransitions {
				if transition.char == ch {
					nextStates[transition.dst] = struct{}{}
				}
			}
		}

		clear(currentStates)
		for nextState := range nextStates {
			appendSetToFirst(currentStates, nextState.buildEpsilonClosure())
			currentStates[nextState] = struct{}{}
		}

		if len(currentStates) == 0 {
			return false, nil
		}
	}
	ok := false
	for state := range currentStates {
		if state.nodeType == GroupEnd && groupsIndexs[state.groupName][1] == -1 {
			groupsIndexs[state.groupName][1] = len(str) - 1
		}
		if nfa.end == state {
			ok = true
		}
	}
	for groupName, indx := range groupsIndexs {
		if indx[0] != -1 && indx[1] != -1 && indx[0] <= indx[1] {
			groups[groupName] = str[indx[0] : indx[1]+1]
		}
	}
	return ok, groups

}

func appendSetToFirst(first, second map[*NfaNode]struct{}) {
	for key := range second {
		_, ok := first[key]
		if !ok {
			first[key] = struct{}{}
		}
	}
}

func (nfa *Nfa) ToGraphviz() string {
	var builder strings.Builder
	visited := make(map[*NfaNode]int)
	nodeCounter := 0

	builder.WriteString("digraph NFA {\n")
	builder.WriteString("    rankdir=LR;\n")
	builder.WriteString("    node [shape=circle];\n")
	builder.WriteString("    start [shape=point];\n")

	var traverse func(node *NfaNode)
	traverse = func(node *NfaNode) {
		if _, exists := visited[node]; exists {
			return
		}

		nodeID := nodeCounter
		visited[node] = nodeID
		nodeCounter++

		// Формируем метку узла
		label := fmt.Sprintf("%d", nodeID)
		if node.nodeType != Default {
			switch node.nodeType {
			case GroupStart:
				label = fmt.Sprintf("<g=%s>", node.groupName)
			case GroupEnd:
				label = fmt.Sprintf("</g=%s>", node.groupName)
			}
		}

		// Определяем стиль узла
		nodeAttrs := []string{fmt.Sprintf("label=\"%s\"", label)}

		if node == nfa.end {
			nodeAttrs = append(nodeAttrs, "shape=doublecircle")
		}

		fmt.Fprintf(&builder, "    %d [%s];\n", nodeID, strings.Join(nodeAttrs, ", "))

		// Обрабатываем ε-переходы
		for _, trans := range node.epsilonTransitions {
			traverse(trans.dst)
			dstID := visited[trans.dst]
			fmt.Fprintf(&builder, "    %d -> %d [label=\"ε\"];\n", nodeID, dstID)
		}

		// Обрабатываем символьные переходы
		for _, trans := range node.charTransitions {
			traverse(trans.dst)
			dstID := visited[trans.dst]
			char := string(trans.char)
			switch trans.char {
			case '"':
				char = "\\\""
			case '\\':
				char = "\\\\"
			}
			fmt.Fprintf(&builder, "    %d -> %d [label=\"%s\"];\n",
				nodeID, dstID, char)
		}
	}

	traverse(nfa.start)

	startID := visited[nfa.start]
	fmt.Fprintf(&builder, "    start -> %d;\n", startID)

	builder.WriteString("}\n")
	return builder.String()
}

func (nfa *Nfa) SaveToDot(filename string) error {
	dot := nfa.ToGraphviz()
	return os.WriteFile(filename, []byte(dot), 0644)
}
