include(FetchContent)
FetchContent_Declare(
        Catch2
        GIT_REPOSITORY https://github.com/catchorg/Catch2.git
        GIT_TAG v3.4.0
)

FetchContent_MakeAvailable(Catch2)

add_executable(line_logger_tests
        test_line_logger.cpp
)

target_link_libraries(line_logger_tests PRIVATE
        line_logger
        Catch2::Catch2WithMain
)

add_test(NAME line_logger_tests COMMAND line_logger_test)