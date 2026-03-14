#pragma once


#include <string>

#include "i_line_logger.h"
#include "i_line_checker.h"

struct yy_buffer_state;

extern int yylex();
extern yy_buffer_state* yy_scan_string(const char*);
extern void yy_delete_buffer(yy_buffer_state*);

extern char* var1_;
extern char* var2_;
extern char* var3_;

extern bool valid_;
enum Tokens {
    ID = 1,
    SNUM,
    ASSIGN,
    OPERATION,
    ERROR,
    SUCCESS,
};

extern char* yytext;
extern int yyleng;

namespace LineChecker
{
    class LexLineChecker final : public ILineChecker{
    public:
        bool checkLine(const std::string& input, ILineLogger &logger) override;
    };
}
