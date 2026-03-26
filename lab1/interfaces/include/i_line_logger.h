#pragma once
#include <string>

namespace LineChecker
{
    class ILineLogger {
    public:
        virtual ~ILineLogger() = default;

        virtual void incrementCount(std::string_view variable) = 0;
        [[nodiscard]] virtual int getVariableCount(std::string_view variable) const = 0;
        [[nodiscard]] virtual std::string generateReport() const = 0;
    };
}
