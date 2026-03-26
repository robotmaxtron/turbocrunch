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
	return e.evaluateInternal(expr, false)
}

func (e *Evaluator) EvaluateUpdateAns(expr string) (string, error) {
	return e.evaluateInternal(expr, true)
}

func (e *Evaluator) evaluateInternal(expr string, updateAns bool) (string, error) {
	cexpr := C.CString(expr)
	defer C.free(unsafe.Pointer(cexpr))
	var res *C.char
	if updateAns {
		res = C.evaluator_evaluate_update_ans(e.ptr, cexpr)
	} else {
		res = C.evaluator_evaluate(e.ptr, cexpr)
	}
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

func SetAngleMode(mode byte) {
	C.evaluator_set_angle_mode(C.char(mode))
}

func GetAngleMode() byte {
	return byte(C.evaluator_get_angle_mode())
}

func (e *Evaluator) GetVariable(name string) (string, bool) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	res := C.evaluator_get_variable(e.ptr, cname)
	if res == nil {
		return "", false
	}
	defer C.free(unsafe.Pointer(res))
	return C.GoString(res), true
}

type Constant struct {
	Name     string
	Value    string
	Category string
}

func GetConstants() []Constant {
	count := int(C.evaluator_get_constants_count())
	constants := make([]Constant, 0, count)
	for i := 0; i < count; i++ {
		cname := C.evaluator_get_constant_name(C.int(i))
		cvalue := C.evaluator_get_constant_value(C.int(i))
		ccat := C.evaluator_get_constant_category(C.int(i))

		if cname != nil && cvalue != nil && ccat != nil {
			constants = append(constants, Constant{
				Name:     C.GoString(cname),
				Value:    C.GoString(cvalue),
				Category: C.GoString(ccat),
			})
		}
		if cname != nil {
			C.free(unsafe.Pointer(cname))
		}
		if cvalue != nil {
			C.free(unsafe.Pointer(cvalue))
		}
		if ccat != nil {
			C.free(unsafe.Pointer(ccat))
		}
	}
	return constants
}

func GetUnits() []string {
	count := int(C.evaluator_get_units_count())
	units := make([]string, 0, count)
	for i := 0; i < count; i++ {
		uname := C.evaluator_get_unit_name(C.int(i))
		if uname != nil {
			units = append(units, C.GoString(uname))
			C.free(unsafe.Pointer(uname))
		}
	}
	return units
}

type FunctionInfo struct {
	Identifier  string
	Name        string
	Description string
}

func GetFunctions() []FunctionInfo {
	count := int(C.evaluator_get_functions_count())
	functions := make([]FunctionInfo, 0, count)
	for i := 0; i < count; i++ {
		cid := C.evaluator_get_function_identifier(C.int(i))
		cname := C.evaluator_get_function_name(C.int(i))
		cusage := C.evaluator_get_function_usage(C.int(i))

		if cid != nil {
			info := FunctionInfo{
				Identifier: C.GoString(cid),
			}
			if cname != nil {
				info.Name = C.GoString(cname)
			}
			if cusage != nil {
				info.Description = C.GoString(cusage)
			}
			functions = append(functions, info)
		}

		if cid != nil {
			C.free(unsafe.Pointer(cid))
		}
		if cname != nil {
			C.free(unsafe.Pointer(cname))
		}
		if cusage != nil {
			C.free(unsafe.Pointer(cusage))
		}
	}
	return functions
}

type UserFunction struct {
	Name       string
	Arguments  []string
	Expression string
}

func (e *Evaluator) GetUserFunctions() []UserFunction {
	count := int(C.evaluator_get_user_functions_count(e.ptr))
	functions := make([]UserFunction, 0, count)
	for i := 0; i < count; i++ {
		cname := C.evaluator_get_user_function_name(e.ptr, C.int(i))
		cargs := C.evaluator_get_user_function_args(e.ptr, C.int(i))
		cexpr := C.evaluator_get_user_function_expression(e.ptr, C.int(i))

		if cname != nil {
			fn := UserFunction{
				Name: C.GoString(cname),
			}
			if cargs != nil {
				fn.Arguments = strings.Split(C.GoString(cargs), ";")
				// Filter out empty strings if any
				var filteredArgs []string
				for _, arg := range fn.Arguments {
					if arg != "" {
						filteredArgs = append(filteredArgs, arg)
					}
				}
				fn.Arguments = filteredArgs
			}
			if cexpr != nil {
				fn.Expression = C.GoString(cexpr)
			}
			functions = append(functions, fn)
		}

		if cname != nil {
			C.free(unsafe.Pointer(cname))
		}
		if cargs != nil {
			C.free(unsafe.Pointer(cargs))
		}
		if cexpr != nil {
			C.free(unsafe.Pointer(cexpr))
		}
	}
	return functions
}

func (e *Evaluator) UnsetUserFunction(name string) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.evaluator_unset_user_function(e.ptr, cname)
}

func (e *Evaluator) IsUserFunctionAssign() bool {
	return int(C.evaluator_is_user_function_assign(e.ptr)) == 1
}

func (e *Evaluator) SaveSession(filename string) error {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	if C.evaluator_session_save(e.ptr, cfilename) == 0 {
		return errors.New("failed to save session")
	}
	return nil
}

func (e *Evaluator) LoadSession(filename string) error {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))
	if C.evaluator_session_load(e.ptr, cfilename) == 0 {
		return errors.New("failed to load session")
	}
	return nil
}
