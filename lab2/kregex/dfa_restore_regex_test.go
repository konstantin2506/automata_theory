package kregex

import (
	"testing"
)

func TestRestoreRegex(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "a",
		},
		{
			"abc",
		},
		{
			"abc...",
		},
		{
			"(def)...",
		},
		{
			"def(def)...",
		},
		{
			"(def)...(def)...",
		},
		{
			"(def(def)...)...",
		},
		{
			"def|q|qq|q...",
		},
		{
			"abc|abc|(abc)...",
		},
		{
			"bad?",
		},
		{
			"(aboba)?",
		},
		{
			"a?",
		},
		{
			"q|q?",
		},
		{
			"q|q|q",
		},
		{
			"a|b|c",
		},
		{
			"(ab)|((a)...|(b)...)",
		},
		{
			"eee",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, _ := BuildAst(tt.name)
			dfa := Minimize(NewDfa(tree))
			regex := RestoreRegex(&dfa)
			rtree, _ := BuildAst(regex)
			fromRegex := Minimize(NewDfa(rtree))

			if !AreIsomorphic(&dfa, &fromRegex) {
				t.Errorf("Not isomorphic: '%s', '%s'", tt.name, regex)
			}
		})
	}
}
