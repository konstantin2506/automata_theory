package ast

import (
	"testing"
)

func TestDfa_Search(t *testing.T) {
	tests := []struct {
		name string
		rgex string
		str  string
		want bool
	}{
		{
			name: "single char match",
			rgex: "a",
			str:  "a",
			want: true,
		},
		{
			name: "single char no match",
			rgex: "a",
			str:  "b",
			want: false,
		},
		{
			name: "concat simple",
			rgex: "ab",
			str:  "ab",
			want: true,
		},
		{
			name: "concat no match",
			rgex: "ab",
			str:  "ac",
			want: false,
		},
		{
			name: "concat extra chars",
			rgex: "ab",
			str:  "abc",
			want: false,
		},
		{
			name: "concat substring",
			rgex: "bc",
			str:  "abcd",
			want: false,
		},
		{
			name: "or first branch",
			rgex: "a|b",
			str:  "a",
			want: true,
		},
		{
			name: "or second branch",
			rgex: "a|b",
			str:  "b",
			want: true,
		},
		{
			name: "or no match",
			rgex: "a|b",
			str:  "c",
			want: false,
		},
		{
			name: "kleene zero",
			rgex: "a...",
			str:  "",
			want: true,
		},
		{
			name: "kleene one",
			rgex: "a...",
			str:  "a",
			want: true,
		},
		{
			name: "kleene many",
			rgex: "a...",
			str:  "aaaa",
			want: true,
		},
		{
			name: "kleene no match",
			rgex: "a...",
			str:  "b",
			want: false,
		},
		{
			name: "optional present",
			rgex: "a?b",
			str:  "ab",
			want: true,
		},
		{
			name: "optional absent",
			rgex: "a?b",
			str:  "b",
			want: true,
		},
		{
			name: "optional no match",
			rgex: "a?b",
			str:  "ac",
			want: false,
		},
		{
			name: "complex concat and kleene",
			rgex: "ab...c",
			str:  "ac",
			want: true,
		},
		{
			name: "complex concat and kleene 2",
			rgex: "ab...c",
			str:  "abbbbc",
			want: true,
		},
		{
			name: "complex concat and kleene 3",
			rgex: "ab...c",
			str:  "abc",
			want: true,
		},
		{
			name: "complex concat and kleene no match",
			rgex: "ab...c",
			str:  "abbbbd",
			want: false,
		},
		{
			name: "or with concat",
			rgex: "ab|cd",
			str:  "ab",
			want: true,
		},
		{
			name: "or with concat 2",
			rgex: "ab|cd",
			str:  "cd",
			want: true,
		},
		{
			name: "or with concat no match",
			rgex: "ab|cd",
			str:  "ac",
			want: false,
		},
		{
			name: "repeat exact",
			rgex: "a{3}",
			str:  "aaa",
			want: true,
		},
		{
			name: "repeat too few",
			rgex: "a{3}",
			str:  "aa",
			want: false,
		},
		{
			name: "repeat too many",
			rgex: "a{3}",
			str:  "aaaa",
			want: false,
		},
		{
			name: "repeat with concat",
			rgex: "a{2}b",
			str:  "aab",
			want: true,
		},
		{
			name: "repeat with concat no match",
			rgex: "a{2}b",
			str:  "ab",
			want: false,
		},
		{
			name: "escaped special char",
			rgex: "%(%",
			str:  "(",
			want: true,
		},
		{
			name: "escaped special char 2",
			rgex: "%)%",
			str:  ")",
			want: true,
		},
		{
			name: "nested groups simple",
			rgex: "(ab)",
			str:  "ab",
			want: true,
		},
		{
			name: "named group",
			rgex: "(<name>ab)",
			str:  "ab",
			want: true,
		},
		{
			name: "empty string always matches",
			rgex: "",
			str:  "",
			want: true,
		},
		{
			name: "complex expression",
			rgex: "(a|b)...c?d",
			str:  "aaaacd",
			want: true,
		},
		{
			name: "complex expression 2",
			rgex: "(a|b)...c?d",
			str:  "bd",
			want: true,
		},
		{
			name: "complex expression 3",
			rgex: "(a|b)...c?d",
			str:  "bbbbcd",
			want: true,
		},
		{
			name: "complex expression no match",
			rgex: "(a|b)...c?d",
			str:  "ccccd",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, err := BuildAst(tt.rgex)
			if err != nil {
				t.Fatalf("BuildAst() error = %v", err)
			}

			dfa := NewDfa(tree)
			got := dfa.Search(tt.str)
			if got != tt.want {
				t.Errorf("Search(%q) with regex %q = %v, want %v", tt.str, tt.rgex, got, tt.want)
			}
		})
	}
}
