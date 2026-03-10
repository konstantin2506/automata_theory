#include "SmcLineChecker.h"

using namespace LineChecker;
bool SmcLineChecker::checkLine(const std::string &input, ILineLogger &logger)
{
   // std::cout << "Input: "<<input <<"\ncorrect_ = " << correct<< std::endl;
    if (input.length() == 0)  return false;
    for (auto ch : input) {
        _fsm.checkNext(ch);
        std::cout << "Checked: {" << ch <<"}\ncorrect_ = " << correct << "\ndone : "<< done << std::endl;
        if (correct == false) {
            std::cout << "Reseted" << std::endl;
            _fsm.reset();
            //std::cout << "correct_ = " << correct << "\ndone : "<< done << std::endl;
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