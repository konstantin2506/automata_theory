#pragma once

#include "i_line_logger.h"
#include "i_line_checker.h"

namespace LineChecker
{
    class RegexLineChecker final: public ILineChecker{
    public:
        bool checkLine(const std::string& input, ILineLogger& logger) override;
    };
}
