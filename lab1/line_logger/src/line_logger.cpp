#include "../include/line_logger.h"
#include <sstream>
using namespace LineChecker;


void LineLogger::incrementCount(std::string_view variable)
{
    variables_[variable.data()] += 1;
}

std::string LineLogger::generateReport() const
{
    std::stringstream report;
    for (const auto& [var, count] : variables_) {
        report << var << " - " << std::to_string(count) << "\n";
    }
    return report.str();
}

int LineLogger::getVariableCount(std::string_view variable)
{
    const auto it = variables_.find(variable.data());
    if (it == variables_.end()) {
        return 0;
    }
    return it->second;
}
