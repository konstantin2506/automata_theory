include(FetchContent)
FetchContent_Declare(
        Catch2
        GIT_REPOSITORY https://github.com/catchorg/Catch2.git
        GIT_TAG v3.4.0
)

FetchContent_MakeAvailable(Catch2)

add_executable(line_checker_tests
        test_line_checker.cpp
)

target_link_libraries(line_checker_tests PRIVATE
        line_checker
        Catch2::Catch2WithMain
)

add_test(NAME line_checker_tests COMMAND line_checker_test)