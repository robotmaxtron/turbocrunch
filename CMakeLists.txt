cmake_minimum_required(VERSION 3.16)
project(speedcrunch_terminal)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Find Qt
find_package(Qt6 COMPONENTS Core QUIET)
if(NOT Qt6_FOUND)
    find_package(Qt5 COMPONENTS Core REQUIRED)
    set(QT_LIBRARIES Qt5::Core)
else()
    set(QT_LIBRARIES Qt6::Core)
endif()

set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)

# SpeedCrunch source files (mostly core and math)
set(SC_DIR ${CMAKE_CURRENT_SOURCE_DIR}/SpeedCrunch/src)

set(SC_CORE_SOURCES
    ${SC_DIR}/core/evaluator.cpp
    ${SC_DIR}/core/functions.cpp
    ${SC_DIR}/core/numberformatter.cpp
    ${SC_DIR}/core/settings.cpp
    ${SC_DIR}/core/session.cpp
    ${SC_DIR}/core/sessionhistory.cpp
    ${SC_DIR}/core/variable.cpp
    ${SC_DIR}/core/userfunction.cpp
    ${SC_DIR}/core/opcode.cpp
    ${SC_DIR}/core/constants.cpp
)

set(SC_MATH_SOURCES
    ${SC_DIR}/math/floatcommon.c
    ${SC_DIR}/math/floatconst.c
    ${SC_DIR}/math/floatconvert.c
    ${SC_DIR}/math/floaterf.c
    ${SC_DIR}/math/floatexp.c
    ${SC_DIR}/math/floatgamma.c
    ${SC_DIR}/math/floathmath.c
    ${SC_DIR}/math/floatio.c
    ${SC_DIR}/math/floatipower.c
    ${SC_DIR}/math/floatlog.c
    ${SC_DIR}/math/floatlogic.c
    ${SC_DIR}/math/floatlong.c
    ${SC_DIR}/math/floatnum.c
    ${SC_DIR}/math/floatpower.c
    ${SC_DIR}/math/floatseries.c
    ${SC_DIR}/math/floattrig.c
    ${SC_DIR}/math/floatincgamma.c
    ${SC_DIR}/math/hmath.cpp
    ${SC_DIR}/math/number.c
    ${SC_DIR}/math/cmath.cpp
    ${SC_DIR}/math/cnumberparser.cpp
    ${SC_DIR}/math/quantity.cpp
    ${SC_DIR}/math/rational.cpp
    ${SC_DIR}/math/units.cpp
)

add_executable(speedcrunch_terminal
    main.cpp
    ${SC_CORE_SOURCES}
    ${SC_MATH_SOURCES}
)

target_include_directories(speedcrunch_terminal PRIVATE
    ${SC_DIR}
    ${SC_DIR}/core
    ${SC_DIR}/math
)

target_link_libraries(speedcrunch_terminal ${QT_LIBRARIES})

# SpeedCrunch version definition
target_compile_definitions(speedcrunch_terminal PRIVATE SPEEDCRUNCH_VERSION="terminal-port")
