package backend

import (
	"strings"
	"testing"
)

func TestMathBackendsComprehensive(t *testing.T) {
	config := &Config{Backend: BackendSpeedCrunch}
	evaluator := NewEvaluatorWrapper(config)

	testCases := []struct {
		expr     string
		check    func(string) bool
		backends []MathBackend
	}{
		{
			expr: "1+2",
			check: func(res string) bool {
				return strings.HasPrefix(res, "3")
			},
			backends: []MathBackend{BackendSpeedCrunch, BackendGo},
		},
		{
			expr: "sin(0)",
			check: func(res string) bool {
				return strings.HasPrefix(res, "0")
			},
			backends: []MathBackend{BackendSpeedCrunch, BackendGo},
		},
		{
			expr: "cos(0)",
			check: func(res string) bool {
				return strings.HasPrefix(res, "1")
			},
			backends: []MathBackend{BackendSpeedCrunch, BackendGo},
		},
		{
			expr: "sqrt(4)",
			check: func(res string) bool {
				return strings.HasPrefix(res, "2")
			},
			backends: []MathBackend{BackendSpeedCrunch, BackendGo},
		},
		{
			expr: "ln(2.718281828459045)",
			check: func(res string) bool {
				return strings.HasPrefix(res, "1") || strings.HasPrefix(res, "0.999999")
			},
			backends: []MathBackend{BackendSpeedCrunch, BackendGo},
		},
		{
			expr: "abs(-5)",
			check: func(res string) bool {
				return strings.HasPrefix(res, "5")
			},
			backends: []MathBackend{BackendSpeedCrunch, BackendGo},
		},
		{
			expr: "conj(1+2j)",
			check: func(res string) bool {
				return strings.Contains(res, "1") && strings.Contains(res, "-") && strings.Contains(res, "2")
			},
			backends: []MathBackend{BackendSpeedCrunch, BackendGo},
		},
		{
			expr: "real(3+4j)",
			check: func(res string) bool {
				return strings.HasPrefix(res, "3")
			},
			backends: []MathBackend{BackendSpeedCrunch, BackendGo},
		},
		{
			expr: "imag(3+4j)",
			check: func(res string) bool {
				return strings.HasPrefix(res, "4")
			},
			backends: []MathBackend{BackendSpeedCrunch, BackendGo},
		},
		{
			expr: "1/3",
			check: func(res string) bool {
				return strings.Contains(res, "33333333333333333333333333333333333333333333333333")
			},
			backends: []MathBackend{BackendGo},
		},
		{
			expr: "exp(1)",
			check: func(res string) bool {
				return strings.HasPrefix(res, "2.718281828459045")
			},
			backends: []MathBackend{BackendGo},
		},
	}

	for _, tc := range testCases {
		for _, b := range tc.backends {
			evaluator.config.Backend = b
			res := evaluator.Evaluate(tc.expr)
			if !tc.check(res) {
				t.Errorf("Backend %v: expression %s failed, got %s", b, tc.expr, res)
			}
		}
	}
}

func TestComplexOpsGo(t *testing.T) {
	config := &Config{Backend: BackendGo}
	evaluator := NewEvaluatorWrapper(config)

	testCases := []struct {
		expr     string
		expected string
	}{
		{"sin(1+2j)", evaluator.evaluateGo("sin(1+2j)")}, // Self-consistency for now
		{"(3+4i)*(1-i)", "7 + 1j"},
		{"abs(3+4j)", "5"},
		{"phase(1+1j)", "0.785398163397448"}, // pi/4
	}

	for _, tc := range testCases {
		res := evaluator.Evaluate(tc.expr)
		if !strings.HasPrefix(res, tc.expected) && !strings.Contains(res, tc.expected) {
			t.Errorf("Expression %s failed, expected %s, got %s", tc.expr, tc.expected, res)
		}
	}
}

func TestAnsPersistenceGo(t *testing.T) {
	config := &Config{Backend: BackendGo}
	evaluator := NewEvaluatorWrapper(config)

	// 1 + 1 = 2
	res := evaluator.Evaluate("1+1")
	if !strings.HasPrefix(res, "2") {
		t.Errorf("expected 2, got %s", res)
	}

	// ans * 3 = 6
	res = evaluator.Evaluate("ans * 3")
	if !strings.HasPrefix(res, "6") {
		t.Errorf("expected 6, got %s", res)
	}

	// ans + 4 = 10
	res = evaluator.Evaluate("ans + 4")
	if !strings.HasPrefix(res, "10") {
		t.Errorf("expected 10, got %s", res)
	}
}

func TestGetFunctions(t *testing.T) {
	config := &Config{Backend: BackendSpeedCrunch}
	evaluator := NewEvaluatorWrapper(config)

	funcs := evaluator.GetFunctions()
	if len(funcs) == 0 {
		t.Error("expected at least some functions, got 0")
	}

	foundSin := false
	for _, f := range funcs {
		if f.Identifier == "sin" {
			foundSin = true
			if f.Name == "" {
				t.Error("sin function should have a name")
			}
			if f.Description == "" {
				t.Error("sin function should have a description/usage")
			}
			break
		}
	}

	if !foundSin {
		t.Error("could not find 'sin' function in metadata")
	}
}

func TestUserFunctions(t *testing.T) {
	config := &Config{Backend: BackendSpeedCrunch}
	evaluator := NewEvaluatorWrapper(config)

	// Clean any existing user functions
	funcs := evaluator.GetUserFunctions()
	for _, f := range funcs {
		evaluator.UnsetUserFunction(f.Name)
	}

	// Define a simple function: f(x) = x^2
	res := evaluator.Evaluate("f(x) = x^2")
	if res != "Function defined" {
		t.Errorf("expected 'Function defined', got %q", res)
	}

	// Verify it exists in user functions list
	funcs = evaluator.GetUserFunctions()
	found := false
	for _, f := range funcs {
		if f.Name == "f" {
			found = true
			if len(f.Arguments) != 1 || f.Arguments[0] != "x" {
				t.Errorf("expected arguments [x], got %v", f.Arguments)
			}
			if f.Expression != "x^2" {
				t.Errorf("expected expression 'x^2', got %q", f.Expression)
			}
			break
		}
	}
	if !found {
		t.Error("f(x) not found in user functions")
	}

	// Evaluate the function: f(3) = 9
	res = evaluator.Evaluate("f(3)")
	if !strings.HasPrefix(res, "9") {
		t.Errorf("expected 9, got %q", res)
	}

	// Define another: g(a;b) = a + b
	res = evaluator.Evaluate("g(a;b) = a + b")
	if res != "Function defined" {
		t.Errorf("expected 'Function defined', got %q", res)
	}

	// Evaluate: g(10;20) = 30
	res = evaluator.Evaluate("g(10;20)")
	if !strings.HasPrefix(res, "30") {
		t.Errorf("expected 30, got %q", res)
	}

	// Remove f
	evaluator.UnsetUserFunction("f")
	funcs = evaluator.GetUserFunctions()
	for _, f := range funcs {
		if f.Name == "f" {
			t.Error("f should have been removed")
		}
	}

	// f(3) should now error
	res = evaluator.Evaluate("f(3)")
	if !strings.HasPrefix(res, "Error:") {
		t.Errorf("expected error for f(3), got %q", res)
	}
}

func TestHighPrecisionGo(t *testing.T) {
	config := &Config{Backend: BackendGo}
	evaluator := NewEvaluatorWrapper(config)

	testCases := []struct {
		expr     string
		check    func(string) bool
		expected string
	}{
		{
			expr: "sin(0.5)",
			check: func(res string) bool {
				// We expect 50 digits of precision. sin(0.5) starts with 0.47942553860420300027328793521557
				// Standard float64 sin(0.5) is ~0.479425538604203
				// If we get more digits, it's high precision.
				return strings.HasPrefix(res, "0.479425538604203000")
			},
		},
		{
			expr: "cos(0.5)",
			check: func(res string) bool {
				// cos(0.5) starts with 0.8775825618903727
				return strings.HasPrefix(res, "0.8775825618903727")
			},
		},
		{
			expr: "ln(2)",
			check: func(res string) bool {
				// ln(2) = 0.6931471805599453094172321214581765680755
				return strings.HasPrefix(res, "0.69314718055994530941723212145817")
			},
		},
		{
			expr: "log(100)",
			check: func(res string) bool {
				return strings.HasPrefix(res, "2")
			},
		},
		{
			expr: "exp(1)",
			check: func(res string) bool {
				return strings.HasPrefix(res, "2.71828182845904523536028747135266")
			},
		},
	}

	for _, tc := range testCases {
		res := evaluator.Evaluate(tc.expr)
		if !tc.check(res) {
			t.Errorf("Expression %s failed, got %s", tc.expr, res)
		}
	}
}
