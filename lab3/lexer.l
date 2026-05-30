
%{
package main

import (
	"strconv"
)


type Lexer struct {
	src []byte
	pos int
	
	// Внутренние переменные, которые golex требует для отслеживания токенов
	start int
	end   int
}

func NewLexer(input []byte) *Lexer {
	return &Lexer{src: input, pos: 0}
}

// Метод для golex, возвращающий текущий байт
func (l *Lexer) next() byte {
	if l.pos >= len(l.src) {
		return 0 // Конец ввода (EOF)
	}
	b := l.src[l.pos]
	l.pos++
	return b
}

// Обертка Lex, которую вызывает goyacc
func (l *Lexer) Lex(lval *yySymType) int {
	return l.scan(lval)
}

func (l *Lexer) scan(lval *yySymType) int {
	c := l.next()
%}

%yyc c
%yyn c = l.next()


%%
    // Перед распознаванием токена golex фиксирует начальную позицию
    l.start = l.pos - 1

[ \t\r\n]+
    // Игнорируем пробелы

"+"
	return PLUS
"-"
    return MINUS
"*"
	return MULT
"/"
    return DIV


"true"
    return TRUE
"false"
    return FALSE



"toLogic"
    return TO_LOGIC_CAST
"toDigit"
    return TO_DIGIT_CAST
"size"
    return SIZE_OPERATOR
"resize"
    return RESIZE_OPERATOR

"["
    return SQUARE_BRACE_LEFT
"]"
    return SQUARE_BRACE_RIGHT
"{"
    return GROUP_BRACE_LEFT
"}"
    return GROUP_BRACE_RIGHT
"("
    return DEFAULT_BRACE_LEFT
")"
    return DEFAULT_BRACE_RIGHT

"digit"
    return DIGIT_TYPE
"logic" 
    return LOGIC_TYPE

"="
    return ASSIGNMENT_OPERATOR

"eq"
    return EQUALS
"lt"
    return LOWER_THEN
"lte"
    return LOWER_THEN_OR_EQUALS
"gt"
    return GREATER_THEN
"gte"
    return GREATER_THEN_OR_EQUALS

","
    return COMMA

";"
    return STATEMENT_END

"!"
    return NOT

"&&"
    return AND
    
"most"
    return MOST




"for"
    return FOR

"stop"
    return STOP

"step"
    return STEP




"check"
    return CHECK

"then"
    return THEN

"otherwise"
    return OTHERWISE




"move"
    return MOVE_OPERATOR

"rotate"
    return ROTATE_OPERATOR

"surroundings"
    return SURROUNDINGS 



"routing"
    return FUNCTION_DECL

"perform"
    return FUNCTION_CALL

"return"
    return RETURN


"please"
    return PLEASE

"thank you"
    return THANK_YOU




[a-zA-Z_][a-zA-Z_0-9]*
    l.end = l.pos
    if c != 0 { l.end-- }
    lval.str = string(l.src[l.start:l.end])
    return NAME

[+-]?0[xX][0-9a-fA-F]+
    l.end = l.pos
    if c != 0 { l.end-- }
    val, _ := strconv.ParseInt(string(l.src[l.start:l.end]), 16, 64)
    lval.num = int(val)
    return HEX

[+-]?[a-fA-F][a-fA-F0-9]*
    l.end = l.pos
    if c != 0 { l.end-- }
    val, _ := strconv.ParseInt(string(l.src[l.start:l.end]), 16, 64)
    lval.num = int(val)
    return HEX

[+-]?0[0-7]+
    l.end = l.pos
    if c != 0 { l.end-- }
    val, _ := strconv.ParseInt(string(l.src[l.start:l.end]), 8, 64)
    lval.num = int(val)
    return OCTAL

[+-]?[1-7][0-7]*
    l.end = l.pos
    if c != 0 { l.end-- }
    val, _ := strconv.Atoi(string(l.src[l.start:l.end]))
    lval.num = val
    return DECIMAL

"0"
    lval.num = 0
    return DECIMAL


\0
	return 0 // EOF

.
	// Игнорируем остальные символы
%%

    return 0 
}
