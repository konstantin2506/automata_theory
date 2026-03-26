#include "regex_line_checker.h"


#include <iostream>
#include <regex>

using namespace LineChecker;

bool RegexLineChecker::checkLine(const std::string& input, ILineLogger &logger)
{
    static const std::regex pattern(R"(^(\d+)\s+([a-zA-Z][a-zA-Z0-9]{0,15})\s*=\s*(\-\d+|\d+|[a-zA-Z][a-zA-Z0-9]{0,15})\s*([+\-*/]\s*(\-\d+|\d+|[a-zA-Z][a-zA-Z0-9]{0,15}))?\s*$)",
    std::regex::icase | std::regex::optimize);

    auto begin = std::sregex_iterator(input.begin(), input.end(), pattern);
    auto end = std::sregex_iterator();
    if (begin == end) {
        return false;
    }

    for (auto it = begin; it != end; ++it) {
        const std::smatch& match = *it;
        logger.incrementCount(match[2].str());
        if (match[4].length() != 0 && isalpha(match[5].str()[0])) {
            logger.incrementCount(match[5].str());
        }
        if (isalpha(match[3].str()[0]) ) {
            logger.incrementCount(match[3].str());
        }
    }
    return true;
}
