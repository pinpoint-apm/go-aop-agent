cmake_minimum_required(VERSION 3.0.0)
project(utest)

set(CMAKE_C_FLAGS   "${CMAKE_C_FLAGS} -g3 -O0 -W -Wall -Wextra -Werror -fprofile-arcs -ftest-coverage")
include (CTest)

add_executable(utest_gox86_asm Args.c  goX86asm.c  Inst.c  table.c)
target_compile_definitions(utest_gox86_asm PUBLIC  -DDEBUG_GOx86_ASM)
target_link_libraries(utest_gox86_asm  rt gcov)
add_executable(utest_pinpoint Args.c  goX86asm.c  Inst.c  pinpoint.c  table.c)
target_compile_definitions(utest_pinpoint PUBLIC  -DTRACE)
target_link_libraries(utest_pinpoint  rt gcov)

add_test(utest_pinpoint_mem utest_pinpoint)
add_test(utest_gox86_asm_mem utest_gox86_asm)
