package ast

import (
	"strings"
	"testing"
)

func TestAstStructure(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic kleene",
			input:    "a...",
			expected: "Kleene a",
		},
		{
			name:     "basic or",
			input:    "(a|c)",
			expected: "Or a c",
		},
		{
			name:     "basic optional",
			input:    "a?",
			expected: "Optional a",
		},
		{
			name:     "basic concat",
			input:    "abcd",
			expected: "Concat a b c d",
		},
		{
			name:     "basic named_group",
			input:    "(<AAA>a)",
			expected: "Named_group[AAA] a",
		},
		{
			name:     "basic repeat",
			input:    "a{666}",
			expected: "Repeat:666 a",
		},
		{
			name:     "empty",
			input:    "",
			expected: "FAIL",
		},
		{
			name:     "empty named_group",
			input:    "(<abc>)",
			expected: "FAIL",
		},
		{
			name:     "named_group of empty",
			input:    "(<abc>())",
			expected: "FAIL",
		},
		{
			name:     "a?b",
			input:    "a?b",
			expected: "Concat Optional a b",
		},
		{
			name:     "ab(b|c)b",
			input:    "ab(b|c)b",
			expected: "Concat a b Or b c b",
		},
		{
			name:     "a|b...",
			input:    "a|b...",
			expected: "Or a Kleene b",
		},
		{
			name:     "(a|b)...",
			input:    "(a|b)...",
			expected: "Kleene Or a b",
		},
		{
			name:     "a|a?|(abc)...|a{2}",
			input:    "a|a?|(abc)...|a{2}",
			expected: "Or a Optional a Kleene Concat a b c Repeat:2 a",
		},
		{
			name:     "(<C>(<A>a)|(<B>b))",
			input:    "(<C>(<A>a)|(<B>b))",
			expected: "Named_group[C] Or Named_group[A] a Named_group[B] b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := BuildAst(tt.input)
			if err != nil {
				if tt.expected != "FAIL" {
					t.Errorf("BuildAst() error = %v", err)
				}
			} else if res := strings.TrimSpace(ast.TraverseRLRSpace()); res != tt.expected {
				t.Errorf("expected: '%s'\n\t    got:      '%s'\n", tt.expected, res)
			}
		})
	}
}
