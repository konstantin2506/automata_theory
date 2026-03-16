#include "lex_line_checker.h"

#include <vector>
using namespace LineChecker;

bool LexLineChecker::checkLine(const std::string &input, ILineLogger &logger)
{
    int token;
    bool valid = false;
    bool finished = false;
    std::vector<std::string> vars;

    auto buff = yy_scan_string(input.c_str());

    //auto buff = yy_scan_buffer((char*)&input[0], input.size());
    while ((token = static_cast<Tokens>(yylex()) ) > 0  && !finished) {
        switch (token) {
            case ID:
                vars.emplace_back(yytext);
                break;
            case ERROR:
                valid = false;
                finished = true;
                break;
            case SUCCESS:
                valid = true;
                finished = true;
                break;
            default:
               break;
        }
    }
    yy_delete_buffer(buff);
    if (valid){
        for (const auto& var : vars) {
            logger.incrementCount(var);
        }

    }
    return valid;
}
