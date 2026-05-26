package kregex

import (
	"testing"
)

func TestDifferenceAA_E__EA_E__AE_A(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		{
			"a",
		},
		{
			"abc",
		},
		{
			"a|b",
		},
		{
			"a|a|b",
		},
		{
			"a...",
		},
		{
			"a...|a",
		},
		{
			"ab|(ab)...",
		},
		{
			"aab|a|aaa|aa...a...b",
		},
		{
			"q|q",
		},
		{
			"OOP|Svyatoslav?",
		},
		{
			"mmm...",
		},
		{
			"(aab|a|a...aa)...|aa...(a...b){100}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, _ := BuildAst(tt.name)
			dfa := Minimize(NewDfa(tree))
			dfa = makeCompleted(dfa, dfa.alphabet)
			empty := Difference(dfa, dfa)
			if !empty.isEmptyLanguage() {
				t.Errorf("DifferenceAA(%s, %s) not empty", tt.name, tt.name)
			}

			diffEA := Difference(empty, dfa)
			if !AreIsomorphic(&diffEA, &empty) {
				t.Errorf("DifferenceEA(empty, '%s') not isomorphic to empty", tt.name)
			}
			diffAE := Difference(dfa, empty)
			if !AreIsomorphic(&diffAE, &dfa) {
				t.Errorf("DifferenceAE('%s', empty) not isomorphic to A", tt.name)
			}
		})
	}
}
