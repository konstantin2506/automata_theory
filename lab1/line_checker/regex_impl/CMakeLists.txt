add_library(regex_line_checker
        src/line_checker.cpp
)

target_include_directories(regex_line_checker PUBLIC
        include

)

target_compile_features(regex_line_checker PUBLIC cxx_std_20)
target_link_libraries(regex_line_checker INTERFACE interfaces)
