find_package(FLEX REQUIRED)

add_library(lex_line_checker
        src/test.cpp
        src/lex_line_checker.cpp
)

target_include_directories(lex_line_checker PUBLIC
        ${CMAKE_CURRENT_SOURCE_DIR}/include
)

target_include_directories(lex_line_checker PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/src
)

target_link_libraries(lex_line_checker PUBLIC interfaces)
target_compile_features(lex_line_checker PUBLIC cxx_std_20)

if(BUILD_TESTING)
    add_subdirectory(tests)
endif()