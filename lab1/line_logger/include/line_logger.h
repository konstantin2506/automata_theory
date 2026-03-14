#pragma once
#include <string>
#include <unordered_map>

#include "i_line_logger.h"

namespace LineChecker
{
    class LineLogger final : public ILineLogger{
        std::unordered_map<std::string, int> variables_;

    public:
        void incrementCount(std::string_view variable) override;
        [[nodiscard]] std::string generateReport() const override;
        [[nodiscard]] int getVariableCount(std::string_view variable) const override;
    };
}

