package ast

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

type GLNode interface {
	CalcNullable(specMap map[Node]*NodeSpec) bool
	CalcFirst(specMap map[Node]*NodeSpec, charNums map[Node]int) []int
	CalcLast(specMap map[Node]*NodeSpec, charNums map[Node]int) []int
}

type (
	FollowMap   map[int][]int
	FirstMap    map[int][]int
	LastMap     map[int][]int
	NullableMap map[int]bool
)

// compile time interface cast precheck
var (
	_ GLNode = (*CharNode)(nil)
	_ GLNode = (*ConcatNode)(nil)
	_ GLNode = (*OrNode)(nil)
	_ GLNode = (*KleeneNode)(nil)
	_ GLNode = (*OptionalNode)(nil)
	_ GLNode = (*RepeatNode)(nil)
	_ GLNode = (*NamedGroupNode)(nil)
)

type DFAState struct {
	ID          int
	Positions   []int
	Transitions map[byte]int
	Accepting   bool
}

type DFA struct {
	States     []*DFAState
	StartState int
	Chars      []byte
	CharNums   map[Node]int
}

func markCharsDev(tree Ast, chars []byte, charNums map[Node]int) []byte {
	// RLR
	if len(chars) == 0 {
		chars = append(chars, '\n')
	}
	root := tree.GetRoot()
	if root.Type() == Char {
		charNums[root] = len(chars)
		chars = append(chars, root.String()[0])
	}
	for _, child := range root.Children() {
		subtree := NewSubAst(child)
		chars = markCharsDev(subtree, chars, charNums)
	}
	return chars
}

func MarkChars(tree *Ast) ([]byte, map[Node]int) {
	charNums := make(map[Node]int)
	chars := markCharsDev(*tree, []byte{}, charNums)

	// finish of regex marker
	endMarker := &CharNode{char: '#', number: 0}
	charNums[endMarker] = len(chars)
	chars = append(chars, '#')

	newRoot := &ConcatNode{
		children: []Node{tree.GetRoot(), endMarker},
	}
	tree.root = newRoot

	return chars, charNums
}

func ComputeNullable(tree *Ast, specMap map[Node]*NodeSpec, charNums map[Node]int) NullableMap {
	tree.GetRoot().(GLNode).CalcNullable(specMap)

	nullable := NullableMap{}
	for node, nodeSpec := range specMap {
		if node.Type() == Char {
			pos := charNums[node]
			nullable[pos] = nodeSpec.IsNullable
		}
	}
	return nullable
}

func ComputeFirst(tree *Ast, specMap map[Node]*NodeSpec, charNums map[Node]int) FirstMap {
	tree.GetRoot().(GLNode).CalcFirst(specMap, charNums)

	first := FirstMap{}
	for node, nodeSpec := range specMap {
		if node.Type() == Char {
			pos := charNums[node]
			first[pos] = nodeSpec.First
		}
	}
	return first
}

func ComputeLast(tree *Ast, specMap map[Node]*NodeSpec, charNums map[Node]int) LastMap {
	tree.GetRoot().(GLNode).CalcLast(specMap, charNums)
	last := LastMap{}
	for node, nodeSpec := range specMap {
		if node.Type() == Char {
			pos := charNums[node]
			last[pos] = nodeSpec.Last
		}
	}
	return last
}

func ComputeFollow(specMap map[Node]*NodeSpec) FollowMap {
	follow := FollowMap{}

	for node := range specMap {
		switch n := node.(type) {
		case *ConcatNode:
			children := node.Children()
			for i := 0; i < len(children)-1; i++ {
				left := children[i]
				right := children[i+1]
				for _, pos := range specMap[left].Last {
					follow[pos] = lo.Union(follow[pos], specMap[right].First)
				}
			}
		case *KleeneNode:
			child := node.Children()[0]

			for _, pos := range specMap[child].Last {
				follow[pos] = lo.Union(follow[pos], specMap[child].First)
			}
		case *RepeatNode:
			count := int(n.count)
			if count > 1 {
				child := node.Children()[0]
				childSpec := specMap[child]

				// Создаём связи между последовательными повторениями
				for i := 0; i < count-1; i++ {
					for _, pos := range childSpec.Last {
						follow[pos] = lo.Union(follow[pos], childSpec.First)
					}
				}
			}
		}
	}
	return follow
}

func BuildDFA(root Node, follow FollowMap, first FirstMap, last LastMap,
	chars []byte, charNums map[Node]int, specMap map[Node]*NodeSpec,
) *DFA {
	// 1. Получаем начальное состояние из first корня
	rootFirst := specMap[root].First
	// Сортируем позиции для детерминизма
	sortedFirst := make([]int, len(rootFirst))
	copy(sortedFirst, rootFirst)
	sort.Ints(sortedFirst)

	// Проверяем, является ли начальное состояние принимающим
	rootLast := specMap[root].Last
	startState := &DFAState{
		ID:          0,
		Positions:   sortedFirst,
		Transitions: make(map[byte]int),
		Accepting:   isAccepting(sortedFirst, rootLast),
	}

	dfaStates := []*DFAState{startState}
	stateMap := make(map[string]int) // для быстрого поиска существующих состояний
	stateMap[positionsKey(sortedFirst)] = 0

	queue := []*DFAState{startState}
	stateID := 1

	// 2. Построение переходов для каждого состояния
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// Группируем позиции по символам
		posByChar := make(map[byte][]int)
		for _, pos := range current.Positions {
			// Находим узел по позиции
			var nodeForPos Node
			for node, num := range charNums {
				if num == pos {
					nodeForPos = node
					break
				}
			}
			if nodeForPos != nil {
				char := chars[pos]
				posByChar[char] = append(posByChar[char], pos)
			}
		}

		// Для каждого символа создаём переход
		for char, positions := range posByChar {
			// Вычисляем follow set для этого символа
			newPositionsMap := make(map[int]bool)
			for _, pos := range positions {
				if followPositions, exists := follow[pos]; exists {
					for _, followPos := range followPositions {
						newPositionsMap[followPos] = true
					}
				}
			}

			// Если нет новых позиций, пропускаем
			if len(newPositionsMap) == 0 {
				continue
			}

			// Преобразуем в отсортированный слайс
			sortedPositions := make([]int, 0, len(newPositionsMap))
			for pos := range newPositionsMap {
				sortedPositions = append(sortedPositions, pos)
			}
			sort.Ints(sortedPositions)

			// Проверяем, существует ли уже такое состояние
			key := positionsKey(sortedPositions)
			if existingID, exists := stateMap[key]; exists {
				// Переход в существующее состояние
				current.Transitions[char] = existingID
			} else {
				// Создаём новое состояние
				newState := &DFAState{
					ID:          stateID,
					Positions:   sortedPositions,
					Transitions: make(map[byte]int),
					Accepting:   isAccepting(sortedPositions, rootLast),
				}

				stateMap[key] = stateID
				stateID++
				dfaStates = append(dfaStates, newState)
				queue = append(queue, newState)
				current.Transitions[char] = newState.ID
			}
		}
	}

	return &DFA{
		States:     dfaStates,
		StartState: 0,
		Chars:      chars,
		CharNums:   charNums,
	}
}

// Вспомогательная функция для создания ключа состояния
func positionsKey(positions []int) string {
	parts := make([]string, len(positions))
	for i, pos := range positions {
		parts[i] = strconv.Itoa(pos)
	}
	return strings.Join(parts, ",")
}

func isAccepting(positions []int, lastPositions []int) bool {
	// Проверяем, содержит ли состояние финальную позицию
	lastSet := make(map[int]bool)
	for _, pos := range lastPositions {
		lastSet[pos] = true
	}

	for _, pos := range positions {
		if lastSet[pos] {
			return true
		}
	}

	return false
}

func (dfa *DFA) String() string {
	var builder strings.Builder
	builder.WriteString("DFA:\n")

	for _, state := range dfa.States {
		builder.WriteString(fmt.Sprintf("State %d: positions=%v, accepting=%v\n",
			state.ID, state.Positions, state.Accepting))

		for char, targetID := range state.Transitions {
			builder.WriteString(fmt.Sprintf("  --[%c]--> State %d\n", char, targetID))
		}
	}

	return builder.String()
}

func (dfa *DFA) GraphVizString() string {
	var builder strings.Builder

	builder.WriteString("digraph DFA {\n")
	builder.WriteString("  rankdir=LR;\n")
	builder.WriteString("  node [shape=circle];\n\n")

	// Начальная точка
	builder.WriteString("  start [shape=point];\n")
	builder.WriteString(fmt.Sprintf("  start -> S%d;\n\n", dfa.StartState))

	// Состояния
	for _, state := range dfa.States {
		shape := ""
		if state.Accepting {
			shape = ", peripheries=2"
		}
		label := fmt.Sprintf("S%d\\n", state.ID)
		builder.WriteString(fmt.Sprintf("  S%d [label=\"%s\"%s];\n",
			state.ID, label, shape))
	}

	builder.WriteString("\n")

	// Переходы
	for _, state := range dfa.States {
		for char, targetID := range state.Transitions {
			charStr := string(char)
			if char == '"' {
				charStr = "\\\""
			} else if char == '\\' {
				charStr = "\\\\"
			}
			builder.WriteString(fmt.Sprintf("  S%d -> S%d [label=\"%s\"];\n",
				state.ID, targetID, charStr))
		}
	}

	builder.WriteString("}\n")
	return builder.String()
}

func positionsToString(positions []int) string {
	if len(positions) == 0 {
		return "{}"
	}
	strs := make([]string, len(positions))
	for i, pos := range positions {
		strs[i] = strconv.Itoa(pos)
	}
	return "{" + strings.Join(strs, ",") + "}"
}

func (dfa *DFA) SaveGraphViz(filename string) error {
	return os.WriteFile(filename, []byte(dfa.GraphVizString()), 0o644)
}
