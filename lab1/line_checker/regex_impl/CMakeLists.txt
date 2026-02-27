add_library(regex_line_checker
        src/regex_line_checker.cpp
)

target_include_directories(regex_line_checker PUBLIC
        ${CMAKE_CURRENT_SOURCE_DIR}/include
)

target_include_directories(regex_line_checker PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/src
)

target_link_libraries(regex_line_checker PUBLIC interfaces)
target_compile_features(regex_line_checker PUBLIC cxx_std_20)

if(BUILD_TESTING)
    add_subdirectory(tests)
endif()
