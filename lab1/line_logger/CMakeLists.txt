add_library(line_logger
        src/line_logger.cpp
)

target_include_directories(line_logger PUBLIC
        ${CMAKE_CURRENT_SOURCE_DIR}/include
)

target_compile_features(line_logger PUBLIC cxx_std_20)

if(BUILD_TESTING)
    add_subdirectory(tests)
endif()