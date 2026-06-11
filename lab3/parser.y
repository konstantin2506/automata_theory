
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

top:
    statement_list
    {
        fmt.Println("\n[Парсер]: Успешное окончание разбора выражения!")
    }
;


statement_list:
    statement
|   statement_list statement    
;
statement:
    variable_init_statement
;

variable_init_statement:
    left_part_assign arifmetic_expr STATEMENT_END
;

left_part_init:
    variable_type name_list ASSIGNMENT_OPERATOR
|   variable_type name_list dimension ASSIGNMENT_OPERATOR

;

variable_type:
    DIGIT_TYPE
|   LOGIC_TYPE    

name_list:
    NAME

|   name_list COMMA NAME
;

arifmetic_expr:
    term
;
term:
    term PLUS factor
    {
        fmt.Println("term -> term + factor")
    }
|   term MINUS factor
    {
        fmt.Println("term -> term - factor")
    }

|   factor
    {
        fmt.Println("term -> factor")
    }
;

factor:
    factor MULT basic_part_ar
    {
        fmt.Println("factor -> factor * basic_part_ar")
    }
|    factor DIV basic_part_ar
    {
        fmt.Println("factor -> factor / basic_part_ar")
    }

|   basic_part_ar
    {
        fmt.Println("factor -> basic_part_ar")
    }
;

basic_part_ar:
    DECIMAL
    {
        fmt.Printf("basic_part_ar -> DECIMAL (%d)\n", $1)
    }
|   DEFAULT_BRACE_LEFT term DEFAULT_BRACE_RIGHT
    {
        fmt.Println("basic_part_ar -> ( term )")
    }

;

dimension:
    SQUARE_BRACE_LEFT dimension_inner SQUARE_BRACE_RIGHT
;
dimension_inner:
    arifmetic_expr
|   dimension_inner COMMA arifmetic_expr
;    


%%

// Функция вывода ошибок, требуемая для yyLexer
func (l *lexerCtx) Error(s string) {
	fmt.Fprintf(os.Stderr, "Ошибка: %s\n", s)
}

