#include "line_logger.h"

#include <algorithm>
#include <ranges>
#include <sstream>
using namespace LineChecker;


void LineLogger::incrementCount(std::string_view variable)
{
    std::string lowerCased{variable};
    std::ranges::transform(lowerCased.begin(), lowerCased.end(), lowerCased.begin(), ::tolower);
    variables_[lowerCased] += 1;
}

std::string LineLogger::generateReport() const
{
    std::stringstream report;
    for (const auto& [var, count] : variables_) {
        report << var << " - " << std::to_string(count) << "\n";
    }
    return report.str();
}

int LineLogger::getVariableCount(std::string_view variable) const
{
    std::string lowerCased{variable};
    std::ranges::transform(lowerCased.begin(), lowerCased.end(), lowerCased.begin(), ::tolower);
    const auto it = variables_.find(lowerCased);
    if (it == variables_.end()) {
        return 0;
    }
    return it->second;
}
