%{
#include "lex_line_checker.h"
#include <string.h>
#include <stdbool.h>

#define YY_BUF_SIZE 262144  // 256KB

extern char* yytext;
extern int yyleng;

bool valid_ = true;
char* var1_ = NULL;
char* var2_ = NULL;
char* var3_ = NULL;

%}

%option noyywrap
%option nounput
%option noinput

%x START
%x AFTER_VARNAME
%x AFTER_EQUALS
%x AFTER_OPERAND1
%x AFTER_OPERATION
%x AFTER_OPERAND2
%x ERROR

digit       [0-9]
letter      [a-zA-Z]
varname     {letter}({letter}|{digit}){0,15}
number      {digit}+
snumber     -?{digit}+
space       [ \t]+
newline      \n

%%

<INITIAL>{
    {number}    {
        BEGIN(START);
    }
    .  {return ERROR;}
}

<START>{
    {space}        {}
    {varname}   {
        BEGIN(AFTER_VARNAME);
        return ID;
    }
    . {BEGIN(INITIAL); return ERROR;}
    <<EOF>>  { BEGIN(INITIAL); return ERROR; }
}

<AFTER_VARNAME>{
    {space}        { }
    "="         { BEGIN(AFTER_EQUALS); }
    . {BEGIN(INITIAL); return ERROR;}
    <<EOF>>        { BEGIN(INITIAL); return ERROR; }
}

<AFTER_EQUALS>{
    {space}        {  }
    {varname} {
        BEGIN(AFTER_OPERAND1);
        return ID;
    }
    {snumber} {
        BEGIN(AFTER_OPERAND1);
        return SNUM;
    }
    . {BEGIN(INITIAL); return ERROR;}
    <<EOF>>   { BEGIN(INITIAL); return ERROR;}
}

<AFTER_OPERAND1>{
    {space}        {  }
    [\+\-\*\/]  {
        BEGIN(AFTER_OPERATION);
    }
    . {BEGIN(INITIAL); return ERROR;}
    <<EOF>>        {
        BEGIN(INITIAL); return SUCCESS;
    }
}

<AFTER_OPERATION>{
    {space}        { }
    {varname} {
        BEGIN(AFTER_OPERAND2);
        return ID;
    }
    {snumber} {BEGIN(AFTER_OPERAND2);
                return SNUM;
    }
    . {BEGIN(INITIAL); return ERROR;}
    <<EOF>>          { BEGIN(INITIAL); return ERROR;}
}

<AFTER_OPERAND2>{
    {space}        { }
    <<EOF>>        {
        BEGIN(INITIAL);
        return SUCCESS;
    }
    . { BEGIN(INITIAL); return ERROR; }
}
%%




