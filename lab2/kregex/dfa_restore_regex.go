package kregex

import (
	"slices"
	"strings"
)

type RegexTransition struct {
	from int
	to   int
}

func NewRegexTransition(from int, to int) RegexTransition {
	return RegexTransition{
		from, to,
	}
}

func getRegexTransitions(dfa *Dfa) map[RegexTransition]string {
	result := map[RegexTransition]string{}

	result[NewRegexTransition(-1, 0)] = ""

	for from, accepting := range dfa.states {
		transitions := dfa.table[from]
		loops := []string{}
		for char, to := range transitions {
			if to != from {
				prev, ok := result[NewRegexTransition(from, to)]
				if !ok {
					result[NewRegexTransition(from, to)] = string(char)
				} else {
					result[NewRegexTransition(from, to)] = prev + "|" + string(char)
				}
			} else { // сразу все петли через "|"
				loops = append(loops, string(char))
			}
		}
		if len(loops) != 0 {
			result[NewRegexTransition(from, from)] = strings.Join(loops, "|")
		}

		if accepting {
			result[NewRegexTransition(from, -2)] = ""
		}
	}
	return result
}

func getInsOutsLoops(state int, transitions map[RegexTransition]string) ([]int, []int, bool) {
	incomings := []int{}
	outcomings := []int{}
	hasLoops := false
	for transition := range transitions {
		if transition.to == state && transition.from != state {
			incomings = append(incomings, transition.from)
		}
		if transition.from == state && transition.to != state {
			outcomings = append(outcomings, transition.to)
		}
		if transition.to == state && transition.from == state {
			hasLoops = true
		}
	}
	return incomings, outcomings, hasLoops
}

func RestoreRegex(dfa *Dfa) string {
	transitions := getRegexTransitions(dfa)
	builder := strings.Builder{}
	for current := range dfa.StatesCount() {
		ins, outs, hasLoops := getInsOutsLoops(current, transitions)
		/*fmt.Println("current:", current)
		fmt.Println("ins:", ins)
		fmt.Println("outs:", outs)
		fmt.Println("hasLoops:", hasLoops)
		*/
		newTransitions := map[RegexTransition]string{}
		for _, in := range ins {
			for _, out := range outs {
				builder.Reset()

				R1, okR1 := transitions[NewRegexTransition(in, current)]
				R2 := ""
				if hasLoops {
					R2 = transitions[NewRegexTransition(current, current)]
				}
				R3, okR3 := transitions[NewRegexTransition(current, out)]
				R4, okR4 := transitions[NewRegexTransition(in, out)]

				canBeEmpty := (R4 == "" && okR4)
				if canBeEmpty {
					builder.WriteString("(")
				}
				if okR4 && !canBeEmpty {

					builder.WriteString(R4)
					builder.WriteString("|")
				}
				if okR1 {
					builder.WriteString(R1)
				}
				if R2 != "" {
					builder.WriteString("(")
					builder.WriteString(R2)
					builder.WriteString(")...")
				}
				if okR3 {
					builder.WriteString(R3)
				}
				if canBeEmpty {
					builder.WriteString(")?")
				}

				newTransitions[NewRegexTransition(in, out)] = builder.String()
				// fmt.Println(in, out, builder.String())
			}
		}
		for _, in := range ins {
			for _, out := range outs {
				delete(transitions, NewRegexTransition(in, out))
				delete(transitions, NewRegexTransition(current, out))
				delete(transitions, NewRegexTransition(in, current))
				transitions[NewRegexTransition(in, out)] = newTransitions[NewRegexTransition(in, out)]
			}
		}
	}
	regex := transitions[NewRegexTransition(-1, -2)]
	return regex
}

type statePair struct {
	first  int
	second int
}

func AreIsomorphic(first, second *Dfa) bool {
	// if len(first.states) != len(second.states) {
	//		return false
	//	}
	if !slices.Equal(first.alphabet, second.alphabet) {
		return false
	}

	queue := []statePair{{0, 0}}
	mapping := map[int]int{0: 0}
	revMapping := map[int]int{0: 0}

	for len(queue) > 0 {
		pair := queue[0]
		queue = queue[1:]
		s1, s2 := pair.first, pair.second

		if first.states[s1] != second.states[s2] {
			return false
		}

		for _, char := range first.alphabet {
			to1, ok1 := first.table[s1][char]
			to2, ok2 := second.table[s2][char]

			if ok1 != ok2 {
				return false // переход есть в одном, но нет в другом
			}
			if !ok1 {
				continue // перехода нет в обоих
			}

			if mapped, exists := mapping[to1]; exists {
				if mapped != to2 {
					return false
				}
			} else if _, used := revMapping[to2]; used {
				return false
			} else {
				mapping[to1] = to2
				revMapping[to2] = to1
				queue = append(queue, statePair{to1, to2})
			}
		}
	}

	return true
}
