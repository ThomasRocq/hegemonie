cmake_minimum_required(VERSION 3.8)
project(GenMap C)

find_package(PkgConfig)
find_package(Threads)

pkg_check_modules(IGRAPH REQUIRED igraph)

set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -g -fPIC -pipe -Wall -Wextra -std=gnu99")
set(CMAKE_C_FLAGS_DEBUG          "-O0 -fno-inline")
set(CMAKE_C_FLAGS_RELEASE        "-O2")
set(CMAKE_C_FLAGS_RELWITHDEBINFO "-O2 -fno-inline")
set(CMAKE_C_FLAGS_MINSIZEREL     "-Os -s")

include_directories(${IGRAPH_INCLUDE_DIR})
link_directories(${IGRAPH_LIBRARY_DIR})

add_executable(gen-map-test test/main.c)
target_link_libraries(gen-map-test ${IGRAPH_LIBRARIES})
