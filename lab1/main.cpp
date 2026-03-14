#include <fstream>
#include <iostream>

#include "line_logger.h"
#include "lex_line_checker.h"


int main() {

    LineChecker::LineLogger logger;
    LineChecker::LexLineChecker checker;
    int c = 0;
    auto handle = [&](std::istream &file)
    {
        std::string line;
        while (std::getline(file, line)) {
            if (line == "q") {
                return;
            }
            ++c;
            auto res = checker.checkLine(line, logger);
            std::cout << c << ": "<< res << std::endl;
        }
    };

    std::cout << "### Проверка строк на корректность ###" << std::endl;
    std::cout << "Ввод:\n\t1.Построчно из терминальчика \n\t2.Из файла" << std::endl;
    int answer = 0;
    std::cin >> answer;
    switch (answer) {
        default:
            std::cout << "Некорректный ответ" << std::endl;

            break;
        case 1:
            std::cin.ignore();
            handle(std::cin);
            break;
        case 2:
            std::string filename;
            std::cout << "Введите имя файла: " << std::flush;
            std::cin.ignore();
            std::getline(std::cin, filename);

            std::cout << filename << std::endl;
            std::ifstream inputFile(std::string("../test_files/") + filename);
            handle(inputFile);
            break;
    }
    std::cout << logger.generateReport();
}
