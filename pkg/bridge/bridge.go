package bridge

/*
#cgo pkg-config: Qt5Core Qt5Widgets
#cgo CXXFLAGS: -I${SRCDIR}/../../SpeedCrunch/src -I${SRCDIR}/../../SpeedCrunch/src/core -I${SRCDIR}/../../SpeedCrunch/src/math -std=c++17 -fPIC
#cgo LDFLAGS: -L${SRCDIR}/../.. -lbridge -lstdc++
#include "bridge.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"strings"
	"unsafe"
)

type Evaluator struct {
	ptr C.EvaluatorPtr
}

func NewEvaluator() *Evaluator {
	return &Evaluator{ptr: C.evaluator_init()}
}

func (e *Evaluator) Evaluate(expr string) (string, error) {
	cexpr := C.CString(expr)
	defer C.free(unsafe.Pointer(cexpr))
	res := C.evaluator_evaluate(e.ptr, cexpr)
	if res == nil {
		return "", errors.New("failed to evaluate expression: memory allocation error")
	}
	defer C.free(unsafe.Pointer(res))
	s := C.GoString(res)
	if strings.HasPrefix(s, "Error: ") {
		return "", errors.New(strings.TrimPrefix(s, "Error: "))
	}
	return s, nil
}
