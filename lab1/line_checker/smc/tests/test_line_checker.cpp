#define CATCH_CONFIG_MAIN
#include <iostream>
#include <ranges>
#include <catch2/catch_test_macros.hpp>

#include "line_logger.h"
#include "SmcLineChecker.h"

using namespace LineChecker;

TEST_CASE("SMC")
{
    SECTION("Basic test")
    {
        std::stringstream input("1 x = 5\n"
                      "2 y = abc\n"
                      "3 z = -10\n"
                      "4 var = value\n"
                      "5 counter = 100500\n"
                      "6 negative = -42\n"
                      "7 camelCase = variable\n"
                      "8 underscore_var = val\n"
                      "9 mix123 = 456abc\n"
                      "10 simple = test\n");

        LineLogger logger;
        SmcLineChecker checker;

        std::string line;
        while (std::getline(input, line)) {
            checker.checkLine(line, logger);
        }
        CHECK(logger.getVariableCount("x") == 1);
        CHECK(logger.getVariableCount("y") == 1);
        CHECK(logger.getVariableCount("z") == 1);
        CHECK(logger.getVariableCount("var") == 1);
        CHECK(logger.getVariableCount("counter") == 1);
        CHECK(logger.getVariableCount("negative") == 1);
        CHECK(logger.getVariableCount("camelCase") == 1);
        CHECK(logger.getVariableCount("underscore_var") == 0);
        CHECK(logger.getVariableCount("mix123") == 0);


    }
    SECTION("Cringe scenery")
    {
        class MockLogger final: public ILineLogger {
        public:
            void incrementCount(std::string_view variable) override {};

            [[nodiscard]] int getVariableCount(std::string_view variable) override {return 0;};

            [[nodiscard]] std::string generateReport() const override {return "";};
        };
        SmcLineChecker checker;
        MockLogger logger;

        CHECK(checker.checkLine("", logger) == false);
        CHECK(checker.checkLine("jvnwknvewnvewvnwv vwkvw = nwej nkewvewv", logger) == false);
        CHECK(checker.checkLine("\n", logger) == false);
        CHECK(checker.checkLine("\n\n", logger) == false);
        CHECK(checker.checkLine(" ", logger) == false);
        CHECK(checker.checkLine("1", logger) == false);
        CHECK(checker.checkLine("1 1 = 1 + 1", logger) == false);
        CHECK(checker.checkLine("a 1 = a + 1", logger) == false);
        CHECK(checker.checkLine("100000000000000000000000000000000000 a = 1 + 1", logger) == true);
        CHECK(checker.checkLine("1000000000000000 a = 1 + 1", logger) == true);
        CHECK(checker.checkLine("10 a                                                                     = 1 + 1", logger) == true);
        CHECK(checker.checkLine("10 aaaaabbbbbcccccdd = 1 + 1", logger) == false);
        CHECK(checker.checkLine("10 a = 1 + 1b", logger) == false);
        CHECK(checker.checkLine("-10 a = 1 + 1b", logger) == false);
        CHECK(checker.checkLine("10 a = 1 *+ 1", logger) == false);
        CHECK(checker.checkLine("10 a = 1 +- 1", logger) == false);
        CHECK(checker.checkLine("10 a = 1 + -1", logger) == true);
        CHECK(checker.checkLine("10 a = 1 + -a", logger) == false);
        CHECK(checker.checkLine("10 a = 1 + a-1", logger) == false);
        CHECK(checker.checkLine("10 a = a - -1", logger) == true);
        CHECK(checker.checkLine("10 a = a+-1", logger) == true);
        CHECK(checker.checkLine("10 a = 1--1", logger) == true);

        CHECK(checker.checkLine("1 a = 1 + a-1", logger) == false);  // оператор - слипся с a
        CHECK(checker.checkLine("2 b = 1 - -1", logger) == true);    // два минуса подряд с пробелом (отрицательное число)
        CHECK(checker.checkLine("3 c = 1--1", logger) == true);     // два минуса без пробела
        CHECK(checker.checkLine("4 d = 1 - - 1", logger) == false);  // пробел между - и числом
        CHECK(checker.checkLine("5 e = 1 +- 1", logger) == false);   // +- вместе
        CHECK(checker.checkLine("6 f = 1 *- 1", logger) == false);   // *- вместе
        CHECK(checker.checkLine("7 g = 1 /- 1", logger) == false);   // /- вместе
        CHECK(checker.checkLine("8 h = 1+1", logger) == true);      // нет пробелов вокруг оператора
        CHECK(checker.checkLine("9 i = 1  +  1", logger) == true);   // много пробелов
        CHECK(checker.checkLine("10 j = 1 +1", logger) == true);    // пробел только слева от оператора
        CHECK(checker.checkLine("11 k=1+1", logger) == true);       // совсем без пробелов
        CHECK(checker.checkLine("12 l = 1 + 1  ", logger) == true);  // пробелы в конце

        CHECK(checker.checkLine("13 _m = 1", logger) == false);      // начинается с подчеркивания
        CHECK(checker.checkLine("14 m_ = 1", logger) == false);       // заканчивается подчеркиванием
        CHECK(checker.checkLine("15 m123 = 1", logger) == true);     // цифры в имени
        CHECK(checker.checkLine("16 1m = 1", logger) == false);      // начинается с цифры
        CHECK(checker.checkLine("17 m = 1", logger) == true);        // одна буква
        CHECK(checker.checkLine("18 verylongvariablename123 = 1", logger) == false); // >16 символов
        CHECK(checker.checkLine("19 aBcDeFgH = 1", logger) == true); // разные регистры

        CHECK(checker.checkLine("20 a = 01", logger) == true);       // ведущий ноль
        CHECK(checker.checkLine("21 a = 001", logger) == true);      // несколько ведущих нулей
        CHECK(checker.checkLine("22 a = -0", logger) == true);       // отрицательный ноль
        CHECK(checker.checkLine("23 a = -01", logger) == true);      // отрицательный с ведущим нулем
        CHECK(checker.checkLine("26 a = 0", logger) == true);        // ноль
    }
}