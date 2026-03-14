#include "SmcLineChecker.h"

using namespace LineChecker;
bool SmcLineChecker::checkLine(const std::string &input, ILineLogger &logger)
{
    if (input.length() == 0)  return false;
    for (auto ch : input) {
        _fsm.checkNext(ch);
        if (correct == false) {
            _fsm.reset();
            return false;
        }
    }
    if (var1Size() ) logger.incrementCount(var1);
    if (var2Size() ) logger.incrementCount(var2);
    if (var3Size() ) logger.incrementCount(var3);
    auto res = done;

    _fsm.setState(Cases::LineNumber);
    resetFSM();
    return res;
}