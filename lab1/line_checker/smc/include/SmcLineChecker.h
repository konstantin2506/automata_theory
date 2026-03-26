#pragma once
#include "../src/SmcLineCheckerFsm.h"
#include "i_line_logger.h"
#include "i_line_checker.h"



namespace LineChecker
{
    class SmcLineChecker final: public ILineChecker{
    public:
        bool checkLine(const std::string &input, ILineLogger &logger) override;

    private:
        std::string var1;
        std::string var2;
        std::string var3;
        SmcLineCheckerFsm _fsm;
        bool correct = true;
        bool done = false;

    public:
        SmcLineChecker() : _fsm(*this)
        {
            startFSM();
        }
        void setInCorrect()
        {
            correct = false;

        }
        void setDone()
        {
            done = true;
        }
        void setNotDone()
        {
            done = false;
        }
        void reset() {}
        void resetFSM()
        {
            correct = true;
            done = false;
            var1.clear();
            var2.clear();
            var3.clear();

        }

        void printStart() { std::cout << "Start" << std::endl; }

    private:
        void startFSM()
        {
            _fsm.enterStartState();
        }


    public:
        void var1PushChar(char ch) {var1 += ch;}
        void var2PushChar(char ch) {var2 += ch;}
        void var3PushChar(char ch) {var3 += ch;}

        [[nodiscard]] size_t var1Size() const {return var1.size();}
        [[nodiscard]] size_t var2Size() const {return var2.size();}
        [[nodiscard]] size_t var3Size() const {return var3.size();}

        void checkNext(char ch) {}

    };
}
