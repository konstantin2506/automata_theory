package ast

import (
	"testing"
)

func TestDfa_Search(t *testing.T) {
	tests := []struct {
		name  string
		regex string
		str   string
		want  string
	}{
		{
			name:  "один символ — точное совпадение",
			regex: "a",
			str:   "a",
			want:  "a",
		},
		{
			name:  "один символ — нет совпадения",
			regex: "a",
			str:   "b",
			want:  "",
		},
		{
			name:  "один символ — в середине строки",
			regex: "a",
			str:   "bbabb",
			want:  "a",
		},
		{
			name:  "один символ — первое из нескольких",
			regex: "a",
			str:   "babab",
			want:  "a",
		},

		{
			name:  "конкатенация — точное совпадение",
			regex: "ab",
			str:   "ab",
			want:  "ab",
		},
		{
			name:  "конкатенация — в середине",
			regex: "ab",
			str:   "xxabxx",
			want:  "ab",
		},
		{
			name:  "конкатенация — первое из нескольких",
			regex: "ab",
			str:   "abxxab",
			want:  "ab",
		},
		{
			name:  "конкатенация — нет совпадения",
			regex: "ab",
			str:   "axxb",
			want:  "",
		},
		{
			name:  "конкатенация — с наложением",
			regex: "aa",
			str:   "aaa",
			want:  "aa",
		},
		{
			name:  "конкатенация — три символа",
			regex: "abc",
			str:   "xxabcxx",
			want:  "abc",
		},

		{
			name:  "альтернатива — первая ветвь",
			regex: "a|b",
			str:   "a",
			want:  "a",
		},
		{
			name:  "альтернатива — вторая ветвь",
			regex: "a|b",
			str:   "b",
			want:  "b",
		},
		{
			name:  "альтернатива — первое вхождение из нескольких",
			regex: "a|b",
			str:   "ccbcca",
			want:  "b",
		},
		{
			name:  "альтернатива — нет совпадения",
			regex: "a|b",
			str:   "ccc",
			want:  "",
		},
		{
			name:  "альтернатива — с конкатенацией",
			regex: "ab|cd",
			str:   "xxcdxx",
			want:  "cd",
		},
		{
			name:  "альтернатива — три ветви",
			regex: "a|b|c",
			str:   "xxcxx",
			want:  "c",
		},

		{
			name:  "звезда — ноль повторений (в начале строки)",
			regex: "a...",
			str:   "b",
			want:  "",
		},
		{
			name:  "звезда — одно повторение",
			regex: "a...",
			str:   "a",
			want:  "a",
		},
		{
			name:  "звезда — много повторений (жадная)",
			regex: "a...",
			str:   "aaaab",
			want:  "aaaa",
		},
		{
			name:  "звезда — в середине строки",
			regex: "a...",
			str:   "bbaaaccc",
			want:  "aaa",
		},
		{
			name:  "звезда — первое вхождение",
			regex: "a...",
			str:   "bbaaaaxxaaaa",
			want:  "aaaa",
		},
		{
			name:  "звезда с конкатенацией — жадная внутри",
			regex: "a...b",
			str:   "aaabbb",
			want:  "aaab",
		},
		{
			name:  "звезда с конкатенацией — минимальное",
			regex: "a...b",
			str:   "ab",
			want:  "ab",
		},
		{
			name:  "звезда с конкатенацией — в середине",
			regex: "a...b",
			str:   "xxaabxx",
			want:  "aab",
		},
		{
			name:  "звезда с конкатенацией — несколько b",
			regex: "a...b...",
			str:   "aabbb",
			want:  "aabbb",
		},

		{
			name:  "опциональность — символ есть",
			regex: "a?b",
			str:   "ab",
			want:  "ab",
		},
		{
			name:  "опциональность — символа нет",
			regex: "a?b",
			str:   "b",
			want:  "b",
		},
		{
			name:  "опциональность — в середине строки",
			regex: "a?b",
			str:   "xxbxx",
			want:  "b",
		},
		{
			name:  "опциональность — нет совпадения",
			regex: "a?b",
			str:   "ac",
			want:  "",
		},

		{
			name:  "повтор — ровно 3",
			regex: "a{3}",
			str:   "aaa",
			want:  "aaa",
		},
		{
			name:  "повтор — в середине",
			regex: "a{3}",
			str:   "xxaaaxx",
			want:  "aaa",
		},
		{
			name:  "повтор — недостаточно (ищем 3, есть только 2)",
			regex: "a{3}",
			str:   "aa",
			want:  "",
		},
		{
			name:  "повтор — больше чем нужно (жадно только 3)",
			regex: "a{3}",
			str:   "aaaa",
			want:  "aaa",
		},
		{
			name:  "повтор — с конкатенацией",
			regex: "a{2}b",
			str:   "xaabx",
			want:  "aab",
		},

		{
			name:  "группа — простая",
			regex: "(ab)",
			str:   "ab",
			want:  "ab",
		},
		{
			name:  "группа — с альтернативой",
			regex: "(a|b)c",
			str:   "xacx",
			want:  "ac",
		},
		{
			name:  "группа — со звездой",
			regex: "(ab)...",
			str:   "abababx",
			want:  "ababab",
		},

		{
			name:  "сложное — (a|b)...c",
			regex: "(a|b)...c",
			str:   "aaac",
			want:  "aaac",
		},
		{
			name:  "сложное — (a|b)...c в середине",
			regex: "(a|b)...c",
			str:   "xxabbcxx",
			want:  "abbc",
		},
		{
			name:  "сложное — a...b|c...d",
			regex: "a...b|c...d",
			str:   "xxcccdxx",
			want:  "cccd",
		},
		{
			name:  "сложное — вложенная звезда",
			regex: "(a...b)...",
			str:   "aababx",
			want:  "aabab",
		},
		{
			name:  "сложное — звезда и опциональность",
			regex: "a...b?c",
			str:   "aaac",
			want:  "aaac",
		},
		{
			name:  "сложное — все вместе",
			regex: "(a|b)...c?d...",
			str:   "xbbcddx",
			want:  "bbcdd",
		},
		{
			name:  "сложное — несколько альтернатив со звёздами",
			regex: "a...|b...|c...",
			str:   "xxbbbxx",
			want:  "bbb",
		},

		{
			name:  "пустая строка — всегда совпадение в начале",
			regex: "",
			str:   "abc",
			want:  "",
		},
		{
			name:  "пустая строка — на пустой строке",
			regex: "",
			str:   "",
			want:  "",
		},
		{
			name:  "нет совпадения вообще",
			regex: "xyz",
			str:   "abcdef",
			want:  "",
		},
		{
			name:  "вся строка — совпадение",
			regex: "a...b...",
			str:   "aaabb",
			want:  "aaabb",
		},
		{
			name:  "перекрывающиеся варианты — первое самое левое",
			regex: "aa",
			str:   "aaa",
			want:  "aa",
		},
		{
			name:  "длинная строка — поиск в середине",
			regex: "hello",
			str:   "xxhelloyy",
			want:  "hello",
		},
		{
			name:  "цифры и буквы",
			regex: "123",
			str:   "ab123cd",
			want:  "123",
		},

		{
			name:  "экранированная скобка",
			regex: "(%(%)",
			str:   "ab(cd",
			want:  "(",
		},
		{
			name:  "экранированная звезда",
			regex: "%.%%.%%.%",
			str:   "a...b",
			want:  "...",
		},

		{
			name:  "конкатенация vs звезда — звезда жадно внутрь",
			regex: "a...b",
			str:   "aaabab",
			want:  "aaab",
		},
		{
			name:  "первое вхождение — не самое длинное в строке",
			regex: "a...",
			str:   "aa b aaaa",
			want:  "aa",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, err := BuildAst(tt.regex)
			if err != nil {
				t.Fatalf("BuildAst(%q) error = %v", tt.regex, err)
			}

			dfa := Minimize(NewDfa(tree))
			got := dfa.Search(tt.str)

			if got != tt.want {
				t.Errorf("Search(%q) with regex %q = %q, want %q",
					tt.str, tt.regex, got, tt.want)
			}
		})
	}
}
