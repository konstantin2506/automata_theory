#define CATCH_CONFIG_MAIN
#include <catch2/catch_test_macros.hpp>

#include "../include/line_logger.h"
using namespace LineChecker;
TEST_CASE("Line logger tests")
{
    SECTION("Increment")
    {
        LineLogger logger;
        logger.incrementCount("a");
        logger.incrementCount("a");
        CHECK(logger.generateReport() == "a - 2\n");
    }
    SECTION("Report")
    {
        LineLogger logger;
        CHECK(logger.generateReport() == "");

        logger.incrementCount("a");
        logger.incrementCount("b");
        logger.incrementCount("b");


        CHECK(logger.generateReport() == "b - 2\na - 1\n");
        CHECKED_ELSE(logger.generateReport() == "a - 1\nb - 2\n") {}
    }
}