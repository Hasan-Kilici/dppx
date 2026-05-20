cmake_minimum_required(VERSION 3.16)

project(dppx_native)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

add_library(
    dppx_native
    STATIC

    src/index.cpp
    src/distance.cpp
    src/bridge.cpp
    src/storage.cpp
    src/segment.cpp
    src/wal.cpp
    src/payload.cpp
    src/hnsw.cpp
    src/heap.cpp
    src/vector.cpp
)

target_include_directories(
    dppx_native
    PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}/include
)