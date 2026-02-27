cmake_minimum_required(VERSION 3.16)
project(lab1)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

option(BUILD_TESTING "Build tests" ON)

add_subdirectory(interfaces)
add_subdirectory(line_logger)
add_subdirectory(line_checker/regex_impl)


add_executable(lab1 main.cpp)
target_link_libraries(lab1 PUBLIC regex_line_checker line_logger)
if(BUILD_TESTING)
    enable_testing()
endif()

