
%{
package main

import (
	"fmt"
	"os"
    intr "interpreter"
)
var __rootAST AstNode = nil
func GetRoot() AstNode{
    return __rootAST
}
%}

%union {
	num int
    str string
    boolean bool
    dType intr.VarT
    node intr.AstNode
    nodeList []intr.AstNode
}


%token PLUS
%token MINUS
%token MULT
%token DIV

%token <boolean> TRUE
%token <boolean> FALSE



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

%token <dType> DIGIT_TYPE
%token <dType> LOGIC_TYPE

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

%type <node> statement

%%

top:
    statement_list
    {
        root := intr.NewStatementListNode()
        for _, statement := range $1.([]intr.AstNode){
            root.PushStatement(statement)
        }
        __rootAST = root
        fmt.Println("\n[Парсер]: Успешное окончание разбора выражения!")
    }
;


statement_list:
    statement
    {
        statements := NewStatementListNode()
        statements.PushStatement($1.(intr.AstNode))
        $$ = statements
    }
|   statement_list statement
    {
        statements := $1.(*StatementListNode)
        statements.PushStatement($1.(intr.AstNode))
        $$ = statements
    }
;
statement:
    variable_init_statement
;

variable_init_statement:
    variable_type name_list ASSIGNMENT_OPERATOR arifmetic_expr STATEMENT_END
    {
        $$ = NewAssignNode()
    }
;



variable_type:
    DIGIT_TYPE 
    {
        $$ = $1
    }
|   LOGIC_TYPE 
    {
        $$ = $1
    }

name_list:
    NAME
    {
        $$ = NewNameNode($1)
    }

|   
;

arifmetic_expr:
    term
    {
        $$ = $1
    }
;
term:
    term PLUS factor
    {
        $$ = NewAddNode($1, $3)
    }
|   term MINUS factor
    {
        $$ = NewSubNode($1, $3)
    }

|   factor
    {
        $$ = $1
    }
;

factor:
    factor MULT basic_part_ar
    {
        $$ = NewMulNode($1, $3)

    }
|    factor DIV basic_part_ar
    {
        $$ = NewDivNode($1, $3)
    }

|   basic_part_ar
    {
        $$ = $1
    }
;

basic_part_ar:
    DECIMAL
    {
        $$ = NewIntegerNode($1)
    }
|   DEFAULT_BRACE_LEFT term DEFAULT_BRACE_RIGHT
    {
        $$ = $2
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

