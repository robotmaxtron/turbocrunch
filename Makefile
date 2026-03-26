# Set PKG_CONFIG_PATH if needed
export PKG_CONFIG_PATH := /usr/local/opt/qt@5/lib/pkgconfig:$(PKG_CONFIG_PATH)

QT_FLAGS := $(shell pkg-config --cflags Qt5Core Qt5Widgets)
QT_LIBS := $(shell pkg-config --libs Qt5Core Qt5Widgets)

MOC := $(shell pkg-config --variable=host_bins Qt5Core)/moc
ifeq ($(MOC),/moc)
    # Fallback if pkg-config doesn't provide host_bins
    MOC := moc
endif

SC_DIR := SpeedCrunch/src
SC_CORE_DIR := $(SC_DIR)/core
SC_MATH_DIR := $(SC_DIR)/math

CXX := clang++
CC := clang
CXXFLAGS := -std=c++17 -fPIC -I$(SC_DIR) -I$(SC_CORE_DIR) -I$(SC_MATH_DIR) $(QT_FLAGS) -DSPEEDCRUNCH_VERSION=\"terminal-port\"
CFLAGS := -fPIC -I$(SC_DIR) -I$(SC_MATH_DIR)

# Minimal sources for Evaluator to work
SC_SOURCES := \
    $(SC_CORE_DIR)/evaluator.cpp \
    $(SC_CORE_DIR)/functions.cpp \
    $(SC_CORE_DIR)/numberformatter.cpp \
    $(SC_CORE_DIR)/settings.cpp \
    $(SC_CORE_DIR)/session.cpp \
    $(SC_CORE_DIR)/sessionhistory.cpp \
    $(SC_CORE_DIR)/variable.cpp \
    $(SC_CORE_DIR)/userfunction.cpp \
    $(SC_CORE_DIR)/opcode.cpp \
    $(SC_CORE_DIR)/constants.cpp \
    $(SC_MATH_DIR)/floatcommon.c \
    $(SC_MATH_DIR)/floatconst.c \
    $(SC_MATH_DIR)/floatconvert.c \
    $(SC_MATH_DIR)/floaterf.c \
    $(SC_MATH_DIR)/floatexp.c \
    $(SC_MATH_DIR)/floatgamma.c \
    $(SC_MATH_DIR)/floathmath.c \
    $(SC_MATH_DIR)/floatio.c \
    $(SC_MATH_DIR)/floatipower.c \
    $(SC_MATH_DIR)/floatlog.c \
    $(SC_MATH_DIR)/floatlogic.c \
    $(SC_MATH_DIR)/floatlong.c \
    $(SC_MATH_DIR)/floatnum.c \
    $(SC_MATH_DIR)/floatpower.c \
    $(SC_MATH_DIR)/floatseries.c \
    $(SC_MATH_DIR)/floattrig.c \
    $(SC_MATH_DIR)/floatincgamma.c \
    $(SC_MATH_DIR)/hmath.cpp \
    $(SC_MATH_DIR)/number.c \
    $(SC_MATH_DIR)/cmath.cpp \
    $(SC_MATH_DIR)/cnumberparser.cpp \
    $(SC_MATH_DIR)/rational.cpp \
    $(SC_MATH_DIR)/quantity.cpp \
    $(SC_MATH_DIR)/units.cpp

MOC_SOURCES := moc_evaluator.cpp moc_functions.cpp moc_constants.cpp
MOC_OBJS := moc_evaluator.o moc_functions.o moc_constants.o

BRIDGE_SOURCES := pkg/bridge/bridge.cpp
SC_OBJS := $(patsubst %.cpp,%.o,$(filter %.cpp,$(SC_SOURCES)))
SC_OBJS += $(patsubst %.c,%.o,$(filter %.c,$(SC_SOURCES)))

all: turbocrunch

LIB_NAME := libbridge.a
BRIDGE_OBJS := pkg/bridge/bridge.o

$(LIB_NAME): $(SC_OBJS) $(MOC_OBJS) $(BRIDGE_OBJS)
	ar rcs $@ $^

pkg/bridge/bridge.o: pkg/bridge/bridge.cpp
	$(CXX) $(CXXFLAGS) -c $< -o $@

moc_evaluator.cpp: $(SC_CORE_DIR)/evaluator.h
	$(MOC) $(QT_FLAGS) $< -o $@

moc_functions.cpp: $(SC_CORE_DIR)/functions.h
	$(MOC) $(QT_FLAGS) $< -o $@

moc_constants.cpp: $(SC_CORE_DIR)/constants.h
	$(MOC) $(QT_FLAGS) $< -o $@

moc_units.cpp: $(SC_MATH_DIR)/units.h
	$(MOC) $(QT_FLAGS) $< -o $@

%.o: %.cpp
	$(CXX) $(CXXFLAGS) -c $< -o $@

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

turbocrunch: $(LIB_NAME) cmd/turbocrunch/main.go pkg/backend/math_wrapper.go
	go build -o turbocrunch ./cmd/turbocrunch

test: $(LIB_NAME)
	go test ./pkg/backend ./cmd/turbocrunch

lint:
	go vet ./...

clean:
	rm -f $(SC_OBJS) $(MOC_OBJS) $(MOC_SOURCES) pkg/bridge/bridge.o $(LIB_NAME) turbocrunch
