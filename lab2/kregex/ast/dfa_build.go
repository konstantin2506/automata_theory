package ast

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

type (
	Dfa struct {
		table    []map[byte]int
		states   []bool
		alphabet []byte
	}
)

type (
	positionKey string

	DfaDev struct {
		table map[positionKey]map[byte]*DfaDevState
	}

	DfaDevState struct {
		id        int
		poskey    positionKey
		accepting bool
	}
)

func (dfa *DfaDev) addState(state *DfaDevState) bool {
	if _, ok := dfa.table[state.poskey]; !ok {
		dfa.table[state.poskey] = make(map[byte]*DfaDevState)
		return true
	}
	return false
}

func (dfa *Dfa) addState(accepting bool) {
	dfa.table = append(dfa.table, make(map[byte]int))
	dfa.states = append(dfa.states, accepting)
}

func (dfa *Dfa) StatesCount() int {
	return len(dfa.states)
}

func (dfa *DfaDev) addTransition(fromKey positionKey, to *DfaDevState, char byte) {
	dfa.table[fromKey][char] = to
}

func createPositionKey(positions []int) positionKey {
	slices.Sort(positions)
	strPoses := make([]string, 0, len(positions))
	for _, pos := range positions {
		strPoses = append(strPoses, strconv.Itoa(pos))
	}
	var pkey positionKey = positionKey(strings.Join(strPoses, ","))
	return pkey
}

func NewDfaDevState(id int, positions []int) DfaDevState {
	return DfaDevState{
		id:        id,
		poskey:    createPositionKey(positions),
		accepting: false,
	}
}

func createAlphabet(chars []byte) []byte {
	set := map[byte]struct{}{}
	for i := 1; i < len(chars); i++ {
		if _, ok := set[chars[i]]; !ok {
			set[chars[i]] = struct{}{}
		}
	}
	alphabet := make([]byte, len(set))
	i := 0
	for char := range set {
		if char == '#' {
			alphabet[len(set)-1] = '#'
		} else {
			alphabet[i] = char
			i++
		}
	}

	return alphabet
}

func NewDfa(tree Ast) Dfa {
	if tree.GetRoot() == nil {
		return Dfa{nil, nil, nil}
	}
	specMap := NewNodeSpecMap()
	chars, charNums := MarkChars(&tree)

	alphabet := createAlphabet(chars)

	ComputeNullable(&tree, specMap, charNums)
	ComputeFirst(&tree, specMap, charNums)
	ComputeLast(&tree, specMap, charNums)
	follow := ComputeFollow(specMap)

	dfaDev := DfaDev{table: make(map[positionKey]map[byte]*DfaDevState)}
	firstState := NewDfaDevState(0, specMap[tree.GetRoot()].First)
	if lo.Contains(specMap[tree.GetRoot()].First, len(chars)-1) {
		firstState.accepting = true
	}
	dfaDev.addState(&firstState)

	sortedAlphabet := alphabet[:len(alphabet)-1]
	slices.Sort(sortedAlphabet)

	dfa := Dfa{table: []map[byte]int{}, states: []bool{}, alphabet: sortedAlphabet}
	dfa.addState(firstState.accepting)

	states := map[positionKey]*DfaDevState{}
	states[firstState.poskey] = &firstState

	q := [][]int{}
	q = append(q, specMap[tree.GetRoot()].First)

	id := 1
	for len(q) > 0 {
		currentState := q[0]
		q = q[1:]
		current := states[createPositionKey(currentState)]

		for _, char := range alphabet {
			union := []int{}
			for _, pos := range currentState {
				if chars[pos] == char {
					union = lo.Union(union, follow[pos])
				}
			}
			if len(union) == 0 {
				continue
			}
			newState := NewDfaDevState(id, union)
			if lo.Contains(union, len(chars)-1) { // '#'
				newState.accepting = true
			}
			if ok := dfaDev.addState(&newState); ok {
				dfa.addState(newState.accepting)
				states[newState.poskey] = &newState
				q = append(q, union)
				id++
			}

			dfaDev.addTransition(current.poskey, states[newState.poskey], char)

			dfa.table[current.id][char] = states[newState.poskey].id

		}
	}

	return dfa
}

func (dfa *Dfa) Search(str string) string {
	if dfa.table == nil {
		return ""
	}

	// if accepting := dfa.states[0]; accepting {
	//		return ""
	//	}

	for i := range len(str) {
		current := 0
		jInx := -1
		for j := i; j < len(str); j++ {
			char := str[j]
			next, ok := dfa.table[current][char]
			if !ok {
				break
			}
			current = next
			if accepting := dfa.states[current]; accepting {
				jInx = j
			}
		}
		if jInx != -1 {
			return str[i : jInx+1]
		}
	}

	return ""
}

func (dfa *Dfa) ToGraphviz() string {
	var sb strings.Builder
	sb.WriteString("digraph DFA {\n")
	sb.WriteString("    rankdir=LR;\n\n")
	sb.WriteString("    start [shape=none, label=\"\"];\n\n")

	for id, accepting := range dfa.states {
		shape := "circle"
		if accepting {
			shape = "doublecircle"
		}
		fmt.Fprintf(&sb, "    %d [label=\"S%d\", shape=%s];\n", id, id, shape)
	}
	sb.WriteString("\n    start -> 0;\n\n")

	for fromId, transitions := range dfa.table {
		for char, to := range transitions {
			fmt.Fprintf(&sb, "    %d -> %d [label=\"%c\"];\n", fromId, to, char)
		}
	}

	sb.WriteString("}\n")
	return sb.String()
}

func (dfa *Dfa) SaveToDot(filename string) error {
	dot := dfa.ToGraphviz()
	return os.WriteFile(filename, []byte(dot), 0o644)
}
