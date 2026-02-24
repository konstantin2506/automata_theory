#pragma once

#include <iosfwd>
#include "line_logger.h"

namespace LineChecker
{

    class ILineChecker {
    public:
        virtual ~ILineChecker() = default;

        virtual bool checkLine(std::istream& input, ILineLogger& logger) = 0;
    };
}
