package kregex

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
	charTransition     CharTransition
	nodeType           NfaNodeT
	groupName          string
}

func (node *NfaNode) addEpsilonWith(dst *NfaNode) {
	node.epsilonTransitions = append(node.epsilonTransitions, EpsilonTransition{dst})
}

type Nfa struct {
	start    *NfaNode
	end      *NfaNode
	closures map[*NfaNode]map[*NfaNode]struct{}
}

func newNfaNode() *NfaNode {
	return &NfaNode{
		epsilonTransitions: make([]EpsilonTransition, 0),
		nodeType:           Default,
	}
}

func buildCharNfa(char byte) Nfa {
	first := newNfaNode()
	second := newNfaNode()
	first.charTransition = CharTransition{char, second}
	return Nfa{first, second, nil}
}

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

	return Nfa{first, last, nil}
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
	return Nfa{first, last, nil}
}

func buildOrNfa(nfas []Nfa) Nfa {
	first := newNfaNode()
	last := newNfaNode()
	for _, nfa := range nfas {
		first.addEpsilonWith(nfa.start)
		nfa.end.addEpsilonWith(last)
	}
	return Nfa{first, last, nil}
}

func buildOptionalNfa(nfa Nfa) Nfa {
	first := newNfaNode()
	last := newNfaNode()
	first.addEpsilonWith(nfa.start)
	first.addEpsilonWith(last)
	nfa.end.addEpsilonWith(last)
	return Nfa{first, last, nil}
}

func BuildFromAst(root Node) Nfa {
	if root == nil {
		return Nfa{nil, nil, nil}
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

func (node *NfaNode) buildEpsilonClosure() map[*NfaNode]struct{} {
	closure := map[*NfaNode]struct{}{}
	stack := []*NfaNode{}
	stack = append(stack, node)
	closure[node] = struct{}{}

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

func (nfa *Nfa) precalcEpsilonClosures() map[*NfaNode]map[*NfaNode]struct{} {
	closures := make(map[*NfaNode]map[*NfaNode]struct{})
	nodes := nfa.getAllNodes()

	for _, node := range nodes {
		closures[node] = node.buildEpsilonClosure()
	}

	return closures
}

func (nfa *Nfa) getAllNodes() []*NfaNode {
	nodes := make([]*NfaNode, 0)
	visited := make(map[*NfaNode]bool)

	var traverse func(node *NfaNode)
	traverse = func(node *NfaNode) {
		if node == nil || visited[node] {
			return
		}

		visited[node] = true
		nodes = append(nodes, node)

		for _, transition := range node.epsilonTransitions {
			traverse(transition.dst)
		}

		if transition := node.charTransition; transition.dst != nil {
			traverse(transition.dst)
		}
	}

	traverse(nfa.start)
	return nodes
}

func (nfa *Nfa) searchDev(s string) (string, map[string]string) {
	if nfa.start == nil && nfa.end == nil && s == "" {
		return "", nil
	}

	if nfa.closures == nil {
		nfa.closures = nfa.precalcEpsilonClosures()
	}

	starts := map[string]int{}
	ends := map[string]int{}

	anchor := newNfaNode()
	anchor.addEpsilonWith(nfa.start)
	nfa.closures[anchor] = anchor.buildEpsilonClosure()

	bestStart := -1
	bestEnd := -1

	states := make(map[*NfaNode]int) // node -> начало матча

	for node := range nfa.closures[anchor] {
		states[node] = 0
	}

	for i := 0; i < len(s); i++ {
		c := s[i]
		// fmt.Println("States_", i, states)

		for node := range nfa.closures[anchor] {
			// Не перезаписываем, если уже есть с более ранним start
			if prev, exists := states[node]; !exists || i < prev {
				states[node] = i
				// fmt.Println("new", &node, i)
			}
		}

		nextStates := make(map[*NfaNode]int)
		for node, start := range states {
			if dst := node.charTransition.dst; node.charTransition.char == c && dst != nil {
				if prev, ok := nextStates[dst]; !ok || start < prev {
					nextStates[node.charTransition.dst] = start
				}
				// fmt.Println("nextStates", start)

				for node := range nfa.closures[dst] {
					if node.nodeType == GroupEnd {
						// fmt.Println("END of : ", node.groupName, string(s[i]))
						ends[node.groupName] = i
					}
				}
			}
			if node.nodeType == GroupStart {
				if _, exists := starts[node.groupName]; !exists {
					// fmt.Println("Start of : ", node.groupName, string(s[i]))
					starts[node.groupName] = i
				}
			}

		}

		states = make(map[*NfaNode]int)
		for node, start := range nextStates {
			for target := range nfa.closures[node] {
				if prev, exists := states[target]; !exists || start < prev {
					states[target] = start
					// fmt.Println("epsilon", start)
				}
			}
		}

		if start, ok := states[nfa.end]; ok {
			end := i + 1
			// fmt.Println(start, bestStart, bestEnd)
			if start < bestStart || bestStart == -1 || (start == bestStart && end > bestEnd) {
				bestStart = start
				bestEnd = end
				// fmt.Println("Complete:", s[start:end])
			}
		}
	}

	delete(nfa.closures, anchor)
	// fmt.Println(starts)
	// fmt.Println(ends)

	if bestStart == -1 {
		return "", nil
	}
	groups := make(map[string]string, len(starts))
	for name, end := range ends {
		groups[name] = s[starts[name] : end+1]
	}

	return s[bestStart:bestEnd], groups
}

func (nfa *Nfa) Search(s string) (string, map[string]string) {
	res, g := nfa.searchDev(s)
	if g == nil {
		return "", nil
	}
	if len(g) == 0 {
		return res, g
	}
	return nfa.searchDev(res)
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

		label := fmt.Sprintf("%d", nodeID)
		if node.nodeType != Default {
			switch node.nodeType {
			case GroupStart:
				label = fmt.Sprintf("<g=%s>", node.groupName)
			case GroupEnd:
				label = fmt.Sprintf("</g=%s>", node.groupName)
			}
		}

		nodeAttrs := []string{fmt.Sprintf("label=\"%s\"", label)}

		if node == nfa.end {
			nodeAttrs = append(nodeAttrs, "shape=doublecircle")
		}

		fmt.Fprintf(&builder, "    %d [%s];\n", nodeID, strings.Join(nodeAttrs, ", "))

		for _, trans := range node.epsilonTransitions {
			traverse(trans.dst)
			dstID := visited[trans.dst]
			fmt.Fprintf(&builder, "    %d -> %d [label=\"ε\"];\n", nodeID, dstID)
		}

		if trans := node.charTransition; trans.dst != nil {

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
	return os.WriteFile(filename, []byte(dot), 0o644)
}
