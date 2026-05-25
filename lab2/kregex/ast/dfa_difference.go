package ast

import "github.com/samber/lo"

func Difference(a, b Dfa) Dfa {
	alphabet := lo.Union(a.alphabet, b.alphabet)

	aFull := makeCompleted(a, alphabet)

	product := buildProduct(aFull, Complement(b), alphabet)

	return Minimize(product)
}

type pair struct {
	a int
	b int
}

func buildProduct(a, b Dfa, alphabet []byte) Dfa {
	pairToID := map[pair]int{}

	startPair := pair{0, 0}
	pairToID[startPair] = 0

	table := []map[byte]int{make(map[byte]int)}
	states := []bool{isDifferenceFinal(a, b, startPair)}

	queue := []pair{startPair}
	nextID := 1

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		currentID := pairToID[current]

		for _, sym := range alphabet {
			nextPair := pair{
				a: a.table[current.a][sym],
				b: b.table[current.b][sym],
			}

			if id, exists := pairToID[nextPair]; exists {
				table[currentID][sym] = id
			} else {
				pairToID[nextPair] = nextID

				table = append(table, make(map[byte]int))
				table[currentID][sym] = nextID

				states = append(states, isDifferenceFinal(a, b, nextPair))

				queue = append(queue, nextPair)
				nextID++
			}
		}
	}

	newAlphabet := make([]byte, len(alphabet))
	copy(newAlphabet, alphabet)

	return Dfa{
		table:    table,
		states:   states,
		alphabet: newAlphabet,
	}
}

func isDifferenceFinal(a, b Dfa, p pair) bool {
	return a.states[p.a] && !b.states[p.b]
}

func (dfa *Dfa) isEmptyLanguage() bool {
	result := false
	for _, accepting := range dfa.states {
		result = result || accepting
	}
	return result
}
