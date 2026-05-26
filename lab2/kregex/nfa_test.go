package kregex

import (
	"testing"
)

func TestNfa_Search(t *testing.T) {
	tests := []struct {
		name       string
		regex      string
		str        string
		wantMatch  string
		wantGroups map[string]string
	}{
		{
			name:      "1",
			regex:     "aa...b...c...|q",
			str:       "aaaaabbbccq",
			wantMatch: "aaaaabbbcc",
		},
		{
			name:      "один символ — точное совпадение",
			regex:     "a",
			str:       "a",
			wantMatch: "a",
		},
		{
			name:      "один символ — нет совпадения",
			regex:     "a",
			str:       "b",
			wantMatch: "",
		},
		{
			name:      "один символ — в середине строки",
			regex:     "a",
			str:       "bbabb",
			wantMatch: "a",
		},
		{
			name:      "один символ — первое из нескольких",
			regex:     "a",
			str:       "babab",
			wantMatch: "a",
		},

		{
			name:      "конкатенация — точное совпадение",
			regex:     "ab",
			str:       "ab",
			wantMatch: "ab",
		},
		{
			name:      "конкатенация — в середине",
			regex:     "ab",
			str:       "xxabxx",
			wantMatch: "ab",
		},
		{
			name:      "конкатенация — первое из нескольких",
			regex:     "ab",
			str:       "abxxab",
			wantMatch: "ab",
		},
		{
			name:      "конкатенация — нет совпадения",
			regex:     "ab",
			str:       "axxb",
			wantMatch: "",
		},
		{
			name:      "конкатенация — с наложением",
			regex:     "aa",
			str:       "aaa",
			wantMatch: "aa",
		},
		{
			name:      "конкатенация — три символа",
			regex:     "abc",
			str:       "xxabcxx",
			wantMatch: "abc",
		},

		{
			name:      "альтернатива — первая ветвь",
			regex:     "a|b",
			str:       "a",
			wantMatch: "a",
		},
		{
			name:      "альтернатива — вторая ветвь",
			regex:     "a|b",
			str:       "b",
			wantMatch: "b",
		},
		{
			name:      "альтернатива — первое вхождение",
			regex:     "a|b",
			str:       "ccbcca",
			wantMatch: "b",
		},
		{
			name:      "альтернатива — нет совпадения",
			regex:     "a|b",
			str:       "ccc",
			wantMatch: "",
		},
		{
			name:      "альтернатива — с конкатенацией",
			regex:     "ab|cd",
			str:       "xxcdxx",
			wantMatch: "cd",
		},

		{
			name:      "звезда — одно повторение",
			regex:     "a...",
			str:       "a",
			wantMatch: "a",
		},
		{
			name:      "звезда — много повторений (жадная)",
			regex:     "a...",
			str:       "aaaab",
			wantMatch: "aaaa",
		},
		{
			name:      "звезда — в середине строки",
			regex:     "a...",
			str:       "bbaaaccc",
			wantMatch: "aaa",
		},
		{
			name:      "звезда — ое вхождение",
			regex:     "a...",
			str:       "bbaaaaxxaaaa",
			wantMatch: "aaaa",
		},
		{
			name:      "звезда с конкатенацией — жадная внутри",
			regex:     "a...b",
			str:       "aaabbb",
			wantMatch: "aaab",
		},
		{
			name:      "звезда с конкатенацией — минимальное",
			regex:     "a...b",
			str:       "ab",
			wantMatch: "ab",
		},
		{
			name:      "звезда с конкатенацией — в середине",
			regex:     "a...b",
			str:       "xxaabxx",
			wantMatch: "aab",
		},
		{
			name:      "звезда с конкатенацией — несколько b",
			regex:     "a...b...",
			str:       "aabbb",
			wantMatch: "aabbb",
		},
		{
			name:      "звезда ab...",
			regex:     "ab...",
			str:       "xaabbbbfa",
			wantMatch: "a",
		},

		{
			name:      "опциональность — символ есть",
			regex:     "a?b",
			str:       "ab",
			wantMatch: "ab",
		},
		{
			name:      "опциональность — символа нет",
			regex:     "a?b",
			str:       "b",
			wantMatch: "b",
		},
		{
			name:      "опциональность — в середине строки",
			regex:     "a?b",
			str:       "xxbxx",
			wantMatch: "b",
		},

		{
			name:      "повтор — ровно 3",
			regex:     "a{3}",
			str:       "aaa",
			wantMatch: "aaa",
		},
		{
			name:      "повтор — в середине",
			regex:     "a{3}",
			str:       "xxaaaxx",
			wantMatch: "aaa",
		},
		{
			name:      "повтор — больше чем нужно (жадно только 3)",
			regex:     "a{3}",
			str:       "aaaa",
			wantMatch: "aaa",
		},
		{
			name:      "повтор — с конкатенацией",
			regex:     "a{2}b",
			str:       "xaabx",
			wantMatch: "aab",
		},

		{
			name:      "группа — простая",
			regex:     "(ab)",
			str:       "ab",
			wantMatch: "ab",
		},
		{
			name:      "группа — с альтернативой",
			regex:     "(a|b)c",
			str:       "xacx",
			wantMatch: "ac",
		},

		{
			name:       "именованная группа — с альтернативой",
			regex:      "(<A>a|b)c",
			str:        "xacx",
			wantMatch:  "ac",
			wantGroups: map[string]string{"A": "a"},
		},
		{
			name:       "именованная группа — захват a*",
			regex:      "(<A>a...)b",
			str:        "aaab",
			wantMatch:  "aaab",
			wantGroups: map[string]string{"A": "aaa"},
		},
		{
			name:       "именованная группа — с опциональностью",
			regex:      "(<A>a?)b",
			str:        "ab",
			wantMatch:  "ab",
			wantGroups: map[string]string{"A": "a"},
		},
		{
			name:       "treesh",
			regex:      "(<A>a...)...|(<BC>b...c?)",
			str:        "qqqbbbbbbccaaaaaaaac",
			wantMatch:  "bbbbbbc",
			wantGroups: map[string]string{"BC": "bbbbbbc"},
		},
		{
			name:       "woooow case",
			regex:      "(<A>(<AA>a?|b?))",
			str:        "cccaabbcc",
			wantMatch:  "a",
			wantGroups: map[string]string{"A": "a", "AA": "a"},
		},
		{
			name:      "именованная группа — нет совпадения",
			regex:     "(<A>abc)",
			str:       "xxabxx",
			wantMatch: "",
		},

		{
			name:      "сложное — (a|b)...c",
			regex:     "(a|b)...c",
			str:       "aaac",
			wantMatch: "aaac",
		},
		{
			name:      "сложное — (a|b)...c в середине",
			regex:     "(a|b)...c",
			str:       "xxabbcxx",
			wantMatch: "abbc",
		},
		{
			name:      "сложное — a...b|c...d",
			regex:     "a...b|c...d",
			str:       "xxcccdxx",
			wantMatch: "cccd",
		},
		{
			name:      "сложное — вложенная звезда",
			regex:     "(a...b)...",
			str:       "aababx",
			wantMatch: "aabab",
		},
		{
			name:      "сложное — все вместе",
			regex:     "(a|b)...c?d...",
			str:       "xbbcddx",
			wantMatch: "bbcdd",
		},
		{
			name:       "сложное — с именованной группой",
			regex:      "(<G>(a|b)...)c",
			str:        "xabbcx",
			wantMatch:  "abbc",
			wantGroups: map[string]string{"G": "abb"},
		},

		{
			name:      "пустая строка — на пустой строке",
			regex:     "",
			str:       "",
			wantMatch: "",
		},
		{
			name:      "нет совпадения вообще",
			regex:     "xyz",
			str:       "abcdef",
			wantMatch: "",
		},
		{
			name:      "вся строка — совпадение",
			regex:     "a...b...",
			str:       "aaabb",
			wantMatch: "aaabb",
		},
		{
			name:      "перекрывающиеся варианты — первое самое левое",
			regex:     "aa",
			str:       "aaa",
			wantMatch: "aa",
		},
		{
			name:      "длинная строка — поиск в середине",
			regex:     "hello",
			str:       "xxhelloyy",
			wantMatch: "hello",
		},

		{
			name:      "конкатенация vs звезда — звезда жадно внутрь",
			regex:     "a...b",
			str:       "aaabab",
			wantMatch: "aaab",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, err := BuildAst(tt.regex)
			if err != nil {
				t.Fatalf("BuildAst(%q) error = %v", tt.regex, err)
			}

			nfa := BuildFromAst(tree.GetRoot())
			gotMatch, gotGroups := nfa.Search(tt.str)

			if gotMatch != tt.wantMatch {
				t.Errorf("Search(%q) with regex %q:\n  match = %q, want %q",
					tt.str, tt.regex, gotMatch, tt.wantMatch)
			}

			if tt.wantGroups != nil {
				if gotGroups != nil && len(gotGroups) != len(tt.wantGroups) {
					t.Errorf("Search(%q) with regex %q:\n  groups count = %d, want %d\n  got:  %v\n  want: %v",
						tt.str, tt.regex, len(gotGroups), len(tt.wantGroups), gotGroups, tt.wantGroups)
					return
				}
				for name, wantValue := range tt.wantGroups {
					gotValue, ok := gotGroups[name]
					if !ok {
						t.Errorf("Search(%q) with regex %q:\n  missing group %q\n  got:  %v\n  want: %v",
							tt.str, tt.regex, name, gotGroups, tt.wantGroups)
						continue
					}
					if gotValue != wantValue {
						t.Errorf("Search(%q) with regex %q:\n  group %q = %q, want %q",
							tt.str, tt.regex, name, gotValue, wantValue)
					}
				}
			}
		})
	}
}
