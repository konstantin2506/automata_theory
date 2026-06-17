
%{
package main

import (
	"fmt"
	"os"
    intr "lab3/interpreter"
)
var __rootAST intr.AstNode = nil
var __goodCount = 0

func GetGood() int{
    return __goodCount
}

func GetRoot() intr.AstNode{
    return __rootAST
}
type paramStruct struct{
    names []string
    types []intr.Variable
}
%}

%union {
	num int
    numSlice []int
    str string
    boolean bool
    dType intr.VarT
    node intr.AstNode
    nodeList []intr.AstNode
    params paramStruct
    variable intr.Variable
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

%type <node> statement statement_list statement_with_thank  variable_init_statement expr term factor basic_part_ar print_statement rotate_statement move_statement  if_statement for_statement variable_assign_statement function_decl_statement return_statement
%type <dType> variable_type
%type <nodeList> dimension arguments
%type <params> params_names 
%type <numSlice> dim_int
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
|   statement_list statement_with_thank
    {
        statements := $1.(*intr.StatementListNode)
        statements.PushStatement($2)
        $$ = statements
    }
;
statement_with_thank:
    statement THANK_YOU
    {
        __goodCount++
        $$ = $1
    }
|   PLEASE statement
    {
        __goodCount++
        $$ = $2
    }


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
|   function_decl_statement
    {
        $$ = $1
    }
|   return_statement
    {
        $$ = $1
    }
|   rotate_statement
    {
        $$ = $1
    }
|   move_statement
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
|   NAME dimension ASSIGNMENT_OPERATOR expr STATEMENT_END
    {
        elem := intr.NewArrayElemNode($2, $1)
        $$ = intr.NewArrayAssignNode(elem, $4)
    }
;    
variable_init_statement:
    variable_type NAME ASSIGNMENT_OPERATOR expr STATEMENT_END
    {
        $$ = intr.NewScalarDeclNode($2, $4, $1)
    }
|   variable_type NAME dimension ASSIGNMENT_OPERATOR expr STATEMENT_END
    {
        $$ = intr.NewArrayDeclNode($1, $2, $3, $5)
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
return_statement:
    RETURN expr STATEMENT_END
    {
        $$ = $2
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
function_decl_statement:
    variable_type FUNCTION_DECL NAME DEFAULT_BRACE_LEFT params_names DEFAULT_BRACE_RIGHT GROUP_BRACE_LEFT statement_list GROUP_BRACE_RIGHT
    {
        var res intr.Variable = nil
        if $1 == intr.Int{
            res = intr.NewVariableInt(0)
        }
        if $1 == intr.Bool{
            res = intr.NewVariableBool(false)
        }
        $$ = intr.NewFunctionDeclNode($3, $5.types, $5.names, res, $8)
    }
    
;
rotate_statement:
     ROTATE_OPERATOR factor STATEMENT_END
    {
        $$ = intr.NewRotateNode($2)
    }
;
move_statement:
    MOVE_OPERATOR expr STATEMENT_END
    {
        $$ = intr.NewMoveNode($2)
    }
; 
params_names:
    params_names COMMA variable_type NAME
    {
        names := $1.names
        types := $1.types
        names = append(names, $4)
        if $3 == intr.Int{
            types = append(types, intr.NewVariableInt(0))
        }
        if $3 == intr.Bool{
            types = append(types, intr.NewVariableBool(false))
        }

        $$ = paramStruct{names, types}
    }
|   params_names COMMA variable_type NAME dim_int
    {
        names := $1.names
        types := $1.types
        names = append(names, $4)
        if $3 == intr.Int{
            arr, _ := intr.NewArray($5, intr.NewVariableInt(0))
            types = append(types, arr)
        }
        if $3 == intr.Bool{
            arr, _ := intr.NewArray($5, intr.NewVariableBool(false))
            types = append(types, arr)
        }

        $$ = paramStruct{names, types}
    } 
|   variable_type NAME
    {
        names := []string{}
        types := []intr.Variable{}
        names = append(names, $2)
        if $1 == intr.Int{
            types = append(types, intr.NewVariableInt(0))
        }
        if $1 == intr.Bool{
            types = append(types, intr.NewVariableBool(false))
        }

        $$ = paramStruct{names, types}
    }
|   variable_type NAME dim_int
    {
        names := []string{}
        types := []intr.Variable{}
        names = append(names, $2)
        if $1 == intr.Int{
            arr, _ := intr.NewArray($3, intr.NewVariableInt(0))
            types = append(types, arr)
        }
        if $1 == intr.Bool{
            arr, _ := intr.NewArray($3, intr.NewVariableBool(false))
            types = append(types, arr)
        }

        $$ = paramStruct{names, types}
    }
;
dim_int:
    dim_int SQUARE_BRACE_LEFT DECIMAL SQUARE_BRACE_RIGHT
    {
        vec := $1
        vec = append(vec, $3)
        $$ = vec
    }
|   SQUARE_BRACE_LEFT DECIMAL SQUARE_BRACE_RIGHT
    {
        $$ = []int{$2}
    }
;
dimension:
    dimension SQUARE_BRACE_LEFT expr SQUARE_BRACE_RIGHT
    {
        nodes := append($1, $3)
        $$ = nodes
    }
|   SQUARE_BRACE_LEFT expr SQUARE_BRACE_RIGHT
    {
        nodes := []intr.AstNode{}
        nodes = append(nodes, $2)
        $$ = nodes
    }
;
arguments:
    arguments COMMA expr
    {
        nodes := append($1, $3)
        $$ = nodes
    }
|   expr
    {
        nodes := []intr.AstNode{}
        nodes = append(nodes, $1)
        $$ = nodes
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
|   SIZE_OPERATOR basic_part_ar
    {
        $$ = intr.NewSizeNode($2)
    }
|   MOST basic_part_ar
    {
        $$ = intr.NewMostNode($2)
    }


|   NOT basic_part_ar
    {
        $$ = intr.NewNotNode($2)
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
|   FUNCTION_CALL NAME DEFAULT_BRACE_LEFT arguments DEFAULT_BRACE_RIGHT
    {
        $$ = intr.NewFunctionCallNode($2, $4)
    }
    

|   SURROUNDINGS
    {
        $$ = intr.NewSurrNode()
    }
|   NAME dimension
    {
        $$ = intr.NewArrayElemNode($2, $1)
    }
|   NAME 
    {
        $$ = intr.NewNameNode($1)
    }
;

   


%%

func (l *lexerCtx) Error(s string) {
	fmt.Fprintf(os.Stderr, "Ошибка: %s\n", s)
}

