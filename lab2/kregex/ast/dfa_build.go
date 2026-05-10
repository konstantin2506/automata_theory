package ast

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

type Dfa struct {
	table map[positionKey]map[byte]*DfaState
	first *DfaState
}

type positionKey string

type DfaState struct {
	id        int
	poskey    positionKey
	accepting bool
}

func (dfa *Dfa) addState(state *DfaState) bool {
	if _, ok := dfa.table[state.poskey]; !ok {
		dfa.table[state.poskey] = make(map[byte]*DfaState)
		return true
	}
	return false
}

func (dfa *Dfa) addTransition(fromKey positionKey, to *DfaState, char byte) {
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

func NewDfaState(id int, positions []int) DfaState {
	return DfaState{
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
	alphabet := make([]byte, 0, len(set))
	for char := range set {
		alphabet = append(alphabet, char)
	}
	return alphabet
}

func NewDfa(tree Ast) Dfa {
	if tree.GetRoot() == nil {
		return Dfa{nil, nil}
	}
	specMap := NewNodeSpecMap()
	chars, charNums := MarkChars(&tree)

	alphabet := createAlphabet(chars)

	ComputeNullable(&tree, specMap, charNums)
	ComputeFirst(&tree, specMap, charNums)
	ComputeLast(&tree, specMap, charNums)
	follow := ComputeFollow(specMap)

	dfa := Dfa{table: make(map[positionKey]map[byte]*DfaState)}
	firstState := NewDfaState(0, specMap[tree.GetRoot()].First)
	if lo.Contains(specMap[tree.GetRoot()].First, len(chars)-1) {
		firstState.accepting = true
	}
	dfa.addState(&firstState)
	dfa.first = &firstState

	q := [][]int{}
	q = append(q, specMap[tree.GetRoot()].First)

	id := 2
	for len(q) > 0 {
		currentState := q[0]
		q = q[1:]

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
			newState := NewDfaState(id, union)
			if lo.Contains(union, len(chars)-1) { // '#'
				newState.accepting = true
			}
			if ok := dfa.addState(&newState); ok {
				q = append(q, union)
				id++
			}
			dfa.addTransition(createPositionKey(currentState), &newState, char)
		}
	}

	return dfa
}

func (dfa *Dfa) Search(str string) string {
	if dfa.first == nil && str == "" {
		return ""
	}

	for i := range len(str) {
		builder := strings.Builder{}
		current := dfa.first
		for j := i; j < len(str); j++ {
			char := str[j]
			if next, ok := dfa.table[current.poskey][char]; ok {
				builder.WriteByte(char)
				if current.accepting {
					return builder.String()
				}
				current = next

			}
		}
		if current.accepting {
			return builder.String()
		}
	}

	return ""
}

func (dfa *Dfa) ToGraphviz() string {
	var sb strings.Builder
	sb.WriteString("digraph DFA {\n")
	sb.WriteString("    rankdir=LR;\n\n")
	sb.WriteString("    start [shape=none, label=\"\"];\n\n")

	// Собираем все уникальные состояния из таблицы
	allStates := map[positionKey]*DfaState{}
	for fromKey, transitions := range dfa.table {
		// fromKey может не быть значением ни в одном переходе
		if _, ok := allStates[fromKey]; !ok {
			allStates[fromKey] = nil // временно, уточним ниже
		}
		for _, to := range transitions {
			allStates[to.poskey] = to
		}
	}

	// Для fromKey, которые остались nil, ищем их среди значений
	for fromKey := range dfa.table {
		if allStates[fromKey] == nil {
			for _, transitions := range dfa.table {
				for _, to := range transitions {
					if to.poskey == fromKey {
						allStates[fromKey] = to
						break
					}
				}
			}
		}
	}

	// Выводим узлы
	for key, state := range allStates {
		shape := "circle"
		if state != nil && state.accepting {
			shape = "doublecircle"
		}
		label := ""
		if state != nil {
			label = fmt.Sprintf("S%d", state.id)
		}
		fmt.Fprintf(&sb, "    \"%s\" [label=\"%s\", shape=%s];\n", key, label, shape)
	}

	// Входная стрелка
	fmt.Fprintf(&sb, "\n    start -> \"%s\";\n\n", dfa.first.poskey)

	// Переходы
	for fromKey, transitions := range dfa.table {
		for char, to := range transitions {
			fmt.Fprintf(&sb, "    \"%s\" -> \"%s\" [label=\"%c\"];\n", fromKey, to.poskey, char)
		}
	}

	sb.WriteString("}\n")
	return sb.String()
}

func (dfa *Dfa) SaveToDot(filename string) error {
	dot := dfa.ToGraphviz()
	return os.WriteFile(filename, []byte(dot), 0644)
}
