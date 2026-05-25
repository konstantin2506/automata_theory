package ast

import "maps"

func Complement(dfa Dfa) Dfa {

	full := makeCompleted(dfa, dfa.alphabet)

	inverted := make([]bool, len(full.states))
	for i := range full.states {
		inverted[i] = !full.states[i]
	}

	result := Dfa{
		table:    full.table,
		states:   inverted,
		alphabet: full.alphabet,
	}

	return Minimize(result)
}

func makeCompleted(dfa Dfa, alphabet []byte) Dfa {
	result := Dfa{
		table:    make([]map[byte]int, len(dfa.table)),
		states:   make([]bool, len(dfa.states)),
		alphabet: make([]byte, len(alphabet)),
	}
	copy(result.alphabet, alphabet)
	copy(result.states, dfa.states)
	for i, trans := range dfa.table {
		result.table[i] = make(map[byte]int)
		maps.Copy(result.table[i], trans)
	}

	trapID := len(result.table)
	result.table = append(result.table, make(map[byte]int))
	result.states = append(result.states, false)

	for _, sym := range alphabet {
		result.table[trapID][sym] = trapID
	}

	for i := range result.table {
		for _, sym := range alphabet {
			if _, ok := result.table[i][sym]; !ok {
				result.table[i][sym] = trapID
			}
		}
	}

	return result
}
