#pragma once

namespace LineChecker
{

    class ILineChecker {
    public:
        virtual ~ILineChecker() = default;
        virtual bool checkLine(const std::string& input, ILineLogger &logger) = 0;
    };
}
