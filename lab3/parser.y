
%{
package main

import (
	"fmt"
	"os"
    intr "lab3/interpreter"
)
var __rootAST intr.AstNode = nil
func GetRoot() intr.AstNode{
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

%token PRINT

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

%token EQ
%token LT
%token LTE
%token GT
%token GTE

%token EQ_MAP
%token LT_MAP
%token LTE_MAP
%token GT_MAP
%token GTE_MAP





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

%type <node> statement statement_list variable_init_statement expr term factor basic_part_ar print_statement if_statement for_statement variable_assign_statement
%type <dType> variable_type
%%

top:
    statement_list
    {
        __rootAST = $1
        fmt.Println("\n[Parser]: AstPart=ok")
    }
;


statement_list:
    statement
    {
        statements := intr.NewStatementListNode()
        statements.PushStatement($1)
        $$ = statements
    }
|   statement_list statement
    {
        statements := $1.(*intr.StatementListNode)
        statements.PushStatement($2)
        $$ = statements
    }
;
statement:
    variable_init_statement
    {
        $$ = $1 
    }
|   variable_assign_statement
    {
        $$ = $1
    }    
|   print_statement    
    {
        $$ = $1
    }
|   if_statement
    {
        $$ = $1
    }
|   for_statement
    {
        $$ = $1
    }
    
;
print_statement:
    PRINT DEFAULT_BRACE_LEFT expr DEFAULT_BRACE_RIGHT STATEMENT_END
    {
        $$ = intr.NewPrintNode($3)
    }
;
variable_assign_statement:
    NAME ASSIGNMENT_OPERATOR expr STATEMENT_END
    {
        $$ = intr.NewAssignNode($1, $3)
    }
;    
variable_init_statement:
    variable_type NAME ASSIGNMENT_OPERATOR expr STATEMENT_END
    {
        $$ = intr.NewScalarDeclNode($2, $4, $1)
    }
;
if_statement:
    
    CHECK expr THEN GROUP_BRACE_LEFT statement_list GROUP_BRACE_RIGHT OTHERWISE GROUP_BRACE_LEFT statement_list GROUP_BRACE_RIGHT
    {
        $$ = intr.NewIfStatementNode($2, $5, $9)
    }
|   CHECK expr THEN GROUP_BRACE_LEFT statement_list GROUP_BRACE_RIGHT
    {
        $$ = intr.NewIfStatementNode($2, $5, nil)
    }   
;

for_statement:
    FOR NAME STOP NAME STEP NAME GROUP_BRACE_LEFT statement_list GROUP_BRACE_RIGHT
    {
        $$ = intr.NewForStatementNode($2, $4, $6, $8)
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

;


expr:
    expr EQ
    {
        $$ = intr.NewEqReduceNode($1)
    }
|   expr LT     
    {
        $$ = intr.NewLtReduceNode($1)
    }
|   expr LTE 
    {
        $$ = intr.NewLteReduceNode($1)
    }
|   expr GT 
    {
        $$ = intr.NewGtReduceNode($1)
    }
|   expr GTE 
    {
        $$ = intr.NewGteReduceNode($1)
    }
|   expr EQ_MAP 
    {
        $$ = intr.NewEqMapNode($1)
    }
|   expr LT_MAP 
    {
        $$ = intr.NewLtMapNode($1)
    }
|   expr LTE_MAP 
    {
        $$ = intr.NewLteMapNode($1)
    }
|   expr GT_MAP 
    {
        $$ = intr.NewGtMapNode($1)
    }
|   expr GTE_MAP 
    {
        $$ = intr.NewGteMapNode($1)
    }
|   term
    {
        $$ = $1
    }

;
term:
    term PLUS factor
    {
        $$ = intr.NewAddNode($1, $3)
    }
|   term MINUS factor
    {
        $$ = intr.NewSubNode($1, $3)
    }

|   factor
    {
        $$ = $1
    }
;

factor:
    factor MULT basic_part_ar
    {
        $$ = intr.NewMulNode($1, $3)

    }
|    factor DIV basic_part_ar
    {
        $$ = intr.NewDivNode($1, $3)
    }
|   factor AND basic_part_ar
    {
        $$ = intr.NewAndNode($1, $3)
    }
|   TO_LOGIC_CAST basic_part_ar
    {
        $$ = intr.NewToBooleanCastNode($2)
    }
|   TO_DIGIT_CAST basic_part_ar
    {
        $$ = intr.NewToIntegerCastNode($2)
    }

|   basic_part_ar
    {
        $$ = $1
    }
;

basic_part_ar:
    TRUE
    {
        $$ = intr.NewBooleanNode($1)
    }
|   FALSE
    {
        $$ = intr.NewBooleanNode($1)
    }
|   DECIMAL
    {
        $$ = intr.NewIntegerNode($1)
        
    }
|   DEFAULT_BRACE_LEFT term DEFAULT_BRACE_RIGHT
    {
        $$ = $2
    }
|   NAME 
    {
        $$ = intr.NewNameNode($1)
    }
;

   


%%

// Функция вывода ошибок, требуемая для yyLexer
func (l *lexerCtx) Error(s string) {
	fmt.Fprintf(os.Stderr, "Ошибка: %s\n", s)
}

