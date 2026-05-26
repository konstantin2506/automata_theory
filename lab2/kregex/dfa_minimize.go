package kregex

import (
	"slices"
	"strconv"
	"strings"
)

func Minimize(dfa Dfa) Dfa {
	if len(dfa.table) == 0 {
		return Dfa{nil, nil, nil}
	}

	groups := initialPartition(dfa)

	groups = refine(dfa, groups)

	return buildMinimized(dfa, groups)
}

func initialPartition(dfa Dfa) [][]int {
	accepting := []int{}
	nonAccepting := []int{}

	for state, isFinal := range dfa.states {
		if isFinal {
			accepting = append(accepting, state)
		} else {
			nonAccepting = append(nonAccepting, state)
		}
	}

	groups := [][]int{}
	if len(nonAccepting) > 0 {
		groups = append(groups, nonAccepting)
	}
	if len(accepting) > 0 {
		groups = append(groups, accepting)
	}
	return groups
}

func getGroupIndex(state int, groups [][]int) int {
	for i, group := range groups {
		if slices.Contains(group, state) {
			return i
		}
	}
	return -1
}

func createSignature(state int, dfa Dfa, groups [][]int) string {
	parts := []string{}
	for _, char := range dfa.alphabet {
		if next, ok := dfa.table[state][char]; ok {
			parts = append(parts, strconv.Itoa(getGroupIndex(next, groups)))
		} else {
			parts = append(parts, "-")
		}
	}
	return strings.Join(parts, ",")
}

func splitGroup(group []int, dfa Dfa, groups [][]int) [][]int {
	sigMap := map[string][]int{} // ex: "1,-,2" --> state

	for _, state := range group {
		sig := createSignature(state, dfa, groups)
		sigMap[sig] = append(sigMap[sig], state)
	}

	if len(sigMap) == 1 {
		return [][]int{group}
	}

	result := [][]int{}
	for _, states := range sigMap {
		result = append(result, states)
	}
	return result
}

func refine(dfa Dfa, groups [][]int) [][]int {
	for {
		newGroups := [][]int{}
		changed := false

		for _, group := range groups {
			parts := splitGroup(group, dfa, groups)
			newGroups = append(newGroups, parts...)
			if len(parts) > 1 {
				changed = true
			}
		}

		groups = newGroups
		if !changed {
			break
		}
	}
	return groups
}

func buildMinimized(dfa Dfa, groups [][]int) Dfa {
	startGroupIdx := -1
	for i, group := range groups {
		if slices.Contains(group, 0) {
			startGroupIdx = i
		}
		if startGroupIdx != -1 {
			break
		}
	}

	if startGroupIdx > 0 {
		groups[0], groups[startGroupIdx] = groups[startGroupIdx], groups[0]
	}
	oldToNew := map[int]int{}
	for newState, group := range groups {
		for _, oldState := range group {
			oldToNew[oldState] = newState
		}
	}

	numStates := len(groups)
	newTable := make([]map[byte]int, numStates)
	newStates := make([]bool, numStates)

	for i := range newTable {
		newTable[i] = make(map[byte]int)
	}

	for newState, group := range groups {
		oldState := group[0]
		if dfa.states[oldState] {
			newStates[newState] = true
		}
	}

	for newState, group := range groups {
		rep := group[0]
		for _, char := range dfa.alphabet {
			if next, ok := dfa.table[rep][char]; ok {
				newTable[newState][char] = oldToNew[next]
			}
		}
	}

	newAlphabet := make([]byte, len(dfa.alphabet))
	copy(newAlphabet, dfa.alphabet)

	return Dfa{
		table:    newTable,
		states:   newStates,
		alphabet: newAlphabet,
	}
}
