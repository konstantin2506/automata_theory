
%{
package main

import (
	"fmt"
	"os"
)
%}

%union {
	num int
    str string
}


%token PLUS
%token MINUS
%token MULT
%token DIV

%token TRUE
%token FALSE



%token TO_LOGIC_CAST
%token TO_DIGIT_CAST
%token SIZE_OPERATOR
%token RESIZE_OPERATOR

%token SQUARE_BRACE_LEFT
%token SQUARE_BRACE_RIGHT
%token GROUP_BRACE_LEFT
%token GROUP_BRACE_RIGHT
%token DEFAULT_BRACE_LEFT
%token DEFAULT_BRACE_RIGHT

%token DIGIT_TYPE
%token LOGIC_TYPE

%token ASSIGNMENT_OPERATOR

%token EQUALS
%token LOWER_THEN
%token LOWER_THEN_OR_EQUALS
%token GREATER_THEN
%token GREATER_THEN_OR_EQUALS

%token COMMA
%token STATEMENT_END

%token NOT
%token AND
%token MOST

%token FOR
%token STOP
%token STEP


%token CHECK
%token THEN
%token OTHERWISE


%token MOVE_OPERATOR
%token ROTATE_OPERATOR
%token SURROUNDINGS 


%token FUNCTION_DECL
%token FUNCTION_CALL
%token RETURN

%token PLEASE
%token THANK_YOU


%token <str> NAME

%token <num> HEX;
%token <num> OCTAL;
%token <num> DECIMAL;




%%

arifmetic_expr:
    DECIMAL
    {
        $$ = $1
    }
    arifmetic_expr PLUS DECIMAL
    {
        
    }

%%

// Функция вывода ошибок, требуемая для yyLexer
func (l *Lexer) Error(s string) {
	fmt.Fprintf(os.Stderr, "Ошибка: %s\n", s)
}

