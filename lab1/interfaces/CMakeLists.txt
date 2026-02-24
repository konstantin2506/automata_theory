add_library(interfaces INTERFACE)

target_include_directories(interfaces INTERFACE
        ${CMAKE_CURRENT_SOURCE_DIR}/include
)

target_compile_features(interfaces INTERFACE cxx_std_20)