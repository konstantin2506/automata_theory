cmake_minimum_required(VERSION 3.16)
project(lab1)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)



option(BUILD_TESTING "Build tests" ON)

add_subdirectory(interfaces)
add_subdirectory(line_logger)
add_subdirectory(line_checker)

add_executable(lab1 main.cpp)
target_compile_options(lab1 PRIVATE -Wall -O3)
target_link_libraries(lab1 PUBLIC line_logger lex_line_checker)

add_executable(benchmark perf_test.cpp)
target_compile_options(benchmark PRIVATE -Wall -O3 -Wl,--stack,16777216)

target_link_libraries(benchmark PUBLIC line_logger lex_line_checker regex_line_checker smc_line_checker)
if(BUILD_TESTING)
    enable_testing()
endif()

