package backend

import (
	"math/cmplx"
	"strconv"
	"strings"
	"testing"
)

// TestMathBackends verifies the functionality of both SpeedCrunch and Go backends.
func TestMathBackends(t *testing.T) {
	config := &Config{
		Backend: BackendSpeedCrunch,
	}
	wrapper := NewEvaluatorWrapper(config)

	tests := []struct {
		name    string
		expr    string
		backend MathBackend
		check   func(t *testing.T, res string)
	}{
		{
			"SpeedCrunch Basic",
			"1+2",
			BackendSpeedCrunch,
			func(t *testing.T, res string) {
				if res != "3" {
					t.Errorf("expected 3, got %s", res)
				}
			},
		},
		{
			"SpeedCrunch Complex",
			"(3+4j)*(1-j)",
			BackendSpeedCrunch,
			func(t *testing.T, res string) {
				if !strings.Contains(res, "7+1j") && !strings.Contains(res, "7 + 1j") {
					t.Errorf("expected 7+1j, got %s", res)
				}
			},
		},
		{
			"SpeedCrunch Precision",
			"sin(pi/4)",
			BackendSpeedCrunch,
			func(t *testing.T, res string) {
				if !strings.HasPrefix(res, "0.7071067811865475") {
					t.Errorf("expected high precision, got %s", res)
				}
			},
		},
		{
			"Go Basic",
			"1+2",
			BackendGo,
			func(t *testing.T, res string) {
				if res != "3" {
					t.Errorf("expected 3, got %s", res)
				}
			},
		},
		{
			"Go Complex",
			"(3+4i)*(1-i)",
			BackendGo,
			func(t *testing.T, res string) {
				if !strings.Contains(res, "7+1i") && !strings.Contains(res, "7 + 1i") && !strings.Contains(res, "7+1j") && !strings.Contains(res, "7 + 1j") {
					t.Errorf("expected 7+1i, got %s", res)
				}
			},
		},
		{
			"Go Precision",
			"sin(pi/4)",
			BackendGo,
			func(t *testing.T, res string) {
				val, err := strconv.ParseFloat(res, 64)
				if err != nil {
					t.Errorf("failed to parse result %s: %v", res, err)
					return
				}
				expected := cmplx.Sin(complex(3.14159265358979323846/4, 0))
				if cmplx.Abs(complex(val, 0)-expected) > 1e-15 {
					t.Errorf("expected %v, got %f", expected, val)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Backend = tt.backend
			res := wrapper.Evaluate(tt.expr)
			tt.check(t, res)
		})
	}
}

// TestBackendSwitching verifies that switching between backends works correctly.
func TestBackendSwitching(t *testing.T) {
	config := &Config{
		Backend: BackendSpeedCrunch,
	}
	wrapper := NewEvaluatorWrapper(config)

	res1 := wrapper.Evaluate("sin(pi/4)")
	config.Backend = BackendGo
	res2 := wrapper.Evaluate("sin(pi/4)")

	if res1 == res2 {
		t.Errorf("Switch test FAILED: both backends produced identical string %s (unlikely for different precisions)", res1)
	}
}
