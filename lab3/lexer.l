
%{
package main

import (
	"strconv"
    "bufio"
)

import (
    intr "lab3/interpreter"
)

type lexerCtx struct{
    src *bufio.Reader
    buf []byte
    current byte    
}


func (l *lexerCtx) getc() byte{
   // fmt.Println("sybol is ", string(l.current), l.current)
    if l.current != 0{
        l.buf = append(l.buf, l.current)        
    }    
    l.current = 0
    if b, err := l.src.ReadByte(); err == nil{
        l.current = b 
    }
    return l.current
        
}

func (l *lexerCtx) Lex(lval *yySymType) int {
	return l.yyLex(lval)
}

func (l *lexerCtx) yyLex(yyLval *yySymType) int {
%}


%yyc l.current
%yyn l.current = l.getc()

D -?[0-9]+

%%
    // Перед распознаванием токена golex фиксирует начальную позицию
    l.buf = l.buf[:0]

[ \t\r\n]+
    // Игнорируем пробелы
"+" return PLUS
"-" return MINUS
"*" return MULT
"/" return DIV

"true"
    yyLval.boolean = true
    return TRUE
"false"
    yyLval.boolean = false
    return FALSE

"print"
    return PRINT

"to_logic" return TO_LOGIC_CAST
"to_digit" return TO_DIGIT_CAST
"size" return SIZE_OPERATOR
"resize" return RESIZE_OPERATOR

"[" return SQUARE_BRACE_LEFT
"]" return SQUARE_BRACE_RIGHT
"{" return GROUP_BRACE_LEFT
"}" return GROUP_BRACE_RIGHT
"(" return DEFAULT_BRACE_LEFT
")" return DEFAULT_BRACE_RIGHT

"digit"
    yyLval.dType = intr.Int
    return DIGIT_TYPE
"logic"
    yyLval.dType = intr.Bool
    return LOGIC_TYPE

"=" return ASSIGNMENT_OPERATOR

"eq" return EQ
"lt" return LT
"lte" return LTE
"gt" return GT
"gte" return GTE

"eq[]" return EQ_MAP
"lt[]" return LT_MAP
"lte[]" return LTE_MAP
"gt[]" return GT_MAP
"gte[]" return GTE_MAP

"," return COMMA
";" return STATEMENT_END

"!" return NOT
"&&" return AND
"most" return MOST

"for" return FOR
"stop" return STOP
"step" return STEP


"check" return CHECK
"then" return THEN
"otherwise" return OTHERWISE


"move" return MOVE_OPERATOR
"rotate" return ROTATE_OPERATOR
"surroundings" return SURROUNDINGS 


"routing" return FUNCTION_DECL
"perform" return FUNCTION_CALL
"return" return RETURN

"please" return PLEASE
"thank you" return THANK_YOU


    

{D}
    val, _ := strconv.Atoi(string(l.buf))
    yyLval.num = val
    return DECIMAL


[A-Za-z_][A-Za-z_0-9]*
    val := string(l.buf)
    yyLval.str = val
    return NAME

\0
	return 0 // EOF

.
	return -1
%%
	return 0
}

