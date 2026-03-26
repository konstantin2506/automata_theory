add_library(smc_line_checker
        src/SmcLineChecker.cpp
        src/SmcLineCheckerFsm.cpp
)

target_include_directories(smc_line_checker PUBLIC
        ${CMAKE_CURRENT_SOURCE_DIR}/include
        /home/kostya2506/smc/smc_7_6_0/lib/C++
)

target_include_directories(smc_line_checker PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/src

)

target_link_libraries(smc_line_checker PUBLIC interfaces)
target_compile_features(smc_line_checker PUBLIC cxx_std_20)

if(BUILD_TESTING)
    add_subdirectory(tests)
endif()
