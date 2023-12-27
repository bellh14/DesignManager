cmake_minimum_required(VERSION 3.25)

project(DOE)

set(CMAKE_CXX_STANDARD 17)

add_executable(DOE src/main.cpp)
include_directories(include)
target_include_directories(DOE PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}/include)

