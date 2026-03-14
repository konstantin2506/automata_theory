
#include <chrono>
#include <fstream>
#include <iostream>
#include <vector>





#include <iostream>
#include <fstream>
#include <string>
#include <chrono>
#include <vector>
#include <iomanip>


#include "lex_line_checker.h"
#include "line_logger.h"
#include "line_checker/regex_impl/include/regex_line_checker.h"
#include "line_checker/smc/include/SmcLineChecker.h"


using namespace LineChecker;

int main() {
    auto smc = SmcLineChecker{};
    auto regex = RegexLineChecker{};
    auto lex = LexLineChecker{};
    auto logger = LineLogger{};

    std::vector sizes = {1000, 5000, 10000, 20000, 30000, 40000, 50000, 60000, 70000, 80000, 90000, 100000};
    const int iterations = 10; // количество повторений для каждой строки

    std::vector<long long> smc_times;
    std::vector<long long> regex_times;
    std::vector<long long> lex_times;

    std::cout << "Замер производительности для разных размеров строк\n";
    std::cout << "================================================\n\n";

    for (auto size : sizes) {
        std::string filename = "../generated_files/" + std::to_string(size) + ".txt";
        std::ifstream file(filename);

        if (!file.is_open()) {
            std::cerr << "Ошибка: не удалось открыть файл " << filename << std::endl;
            continue;
        }

        std::string line;
        std::getline(file, line);
        file.close();

        std::cout << "Размер строки: " << size << " символов\n";
        std::cout << "----------------------------------------\n";

        // Замер SMC
        auto startSMC = std::chrono::high_resolution_clock::now();
        for (int i = 0; i < iterations; i++) {
            smc.checkLine(line, logger);
        }
        auto endSMC = std::chrono::high_resolution_clock::now();
        auto smcDuration = std::chrono::duration_cast<std::chrono::microseconds>(endSMC - startSMC);
        long long smcTime = smcDuration.count() / iterations; // среднее время
        smc_times.push_back(smcTime);
        std::cout << "SMC:   " << std::setw(8) << smcTime << " мкс (среднее за " << iterations << " запусков)\n";

        // Замер Regex
        auto startRegex = std::chrono::high_resolution_clock::now();
        for (int i = 0; i < iterations; i++) {
            regex.checkLine(line, logger);
        }
        auto endRegex = std::chrono::high_resolution_clock::now();
        auto regexDuration = std::chrono::duration_cast<std::chrono::microseconds>(endRegex - startRegex);
        long long regexTime = regexDuration.count() / iterations;
        regex_times.push_back(regexTime);
        std::cout << "Regex: " << std::setw(8) << regexTime << " мкс\n";

        // Замер Lex
        auto startLex = std::chrono::high_resolution_clock::now();
        for (int i = 0; i < iterations; i++) {
            lex.checkLine(line, logger);
        }
        auto endLex = std::chrono::high_resolution_clock::now();
        auto lexDuration = std::chrono::duration_cast<std::chrono::microseconds>(endLex - startLex);
        long long lexTime = lexDuration.count() / iterations;
        lex_times.push_back(lexTime);
        std::cout << "Lex:   " << std::setw(8) << lexTime << " мкс\n\n";
    }

    // Вывод результатов в виде таблицы
    std::cout << "\n========== РЕЗУЛЬТАТЫ (среднее время в микросекундах) ==========\n";
    std::cout << std::setw(12) << "Размер"
              << std::setw(12) << "SMC"
              << std::setw(12) << "Regex"
              << std::setw(12) << "Lex" << "\n";
    std::cout << "--------------------------------------------------------\n";

    int i = 0;
    for (auto size: sizes) {
        std::cout << std::setw(12) << size
                  << std::setw(12) << smc_times[i]
                  << std::setw(12) << regex_times[i]
                  << std::setw(12) << lex_times[i] << "\n";
        i++;
    }

    // Сохранение в CSV файл для построения графиков
    std::ofstream csv("results.csv");
    csv << "Size,SMC,Regex,Lex\n";
    for (size_t i = 0; i < sizes.size(); i++) {
        csv << sizes[i] << "," << smc_times[i] << "," << regex_times[i] << "," << lex_times[i] << "\n";
    }
    csv.close();

    std::cout << "\nРезультаты сохранены в results.csv\n";

    return 0;
}