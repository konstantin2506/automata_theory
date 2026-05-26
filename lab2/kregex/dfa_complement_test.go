package kregex

import (
	"testing"
)

func TestComplement(t *testing.T) {
	tests := []struct {
		name string
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
			completedDfa := makeCompleted(dfa, dfa.alphabet)
			doubleComp := Complement(Complement(dfa))
			if !AreIsomorphic(&completedDfa, &doubleComp) {
				t.Errorf("Complement(Complement(dfa(%s)) isnt isomorphic", tt.name)
			}
		})
	}
}
