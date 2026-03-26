// Package main provides a build verification tool for the TurboCrunch project.
package main

import (
	"fmt"
	"math/cmplx"
	"os"
	"strconv"
	"strings"
)

func main() {
	config := &Config{
		Backend: BackendSpeedCrunch,
	}
	wrapper := NewEvaluatorWrapper(config)

	tests := []struct {
		name    string
		expr    string
		backend MathBackend
		check   func(res string) error
	}{
		{
			"SpeedCrunch Basic",
			"1+2",
			BackendSpeedCrunch,
			func(res string) error {
				if res != "3" {
					return fmt.Errorf("expected 3, got %s", res)
				}
				return nil
			},
		},
		{
			"SpeedCrunch Complex",
			"(3+4j)*(1-j)",
			BackendSpeedCrunch,
			func(res string) error {
				if !strings.Contains(res, "7+1j") && !strings.Contains(res, "7 + 1j") {
					return fmt.Errorf("expected 7+1j, got %s", res)
				}
				return nil
			},
		},
		{
			"SpeedCrunch Precision",
			"sin(pi/4)",
			BackendSpeedCrunch,
			func(res string) error {
				if !strings.HasPrefix(res, "0.7071067811865475") {
					return fmt.Errorf("expected high precision, got %s", res)
				}
				return nil
			},
		},
		{
			"Go Basic",
			"1+2",
			BackendGo,
			func(res string) error {
				if res != "3" {
					return fmt.Errorf("expected 3, got %s", res)
				}
				return nil
			},
		},
		{
			"Go Complex",
			"(3+4i)*(1-i)",
			BackendGo,
			func(res string) error {
				if !strings.Contains(res, "7+1i") && !strings.Contains(res, "7 + 1i") {
					return fmt.Errorf("expected 7+1i, got %s", res)
				}
				return nil
			},
		},
		{
			"Go Precision",
			"sin(pi/4)",
			BackendGo,
			func(res string) error {
				val, err := strconv.ParseFloat(res, 64)
				if err != nil {
					return fmt.Errorf("failed to parse result %s: %v", res, err)
				}
				expected := cmplx.Sin(complex(3.14159265358979323846/4, 0))
				if cmplx.Abs(complex(val, 0)-expected) > 1e-15 {
					return fmt.Errorf("expected %v, got %f", expected, val)
				}
				return nil
			},
		},
	}

	failed := 0
	for i, tt := range tests {
		config.Backend = tt.backend
		res := wrapper.Evaluate(tt.expr)
		if err := tt.check(res); err != nil {
			fmt.Printf("Test %d FAILED [%s] (%s, %v): %v\n", i, tt.name, tt.expr, tt.backend, err)
			failed++
		} else {
			fmt.Printf("Test %d PASSED [%s] (%s, %v): %s\n", i, tt.name, tt.expr, tt.backend, res)
		}
	}

	// Test Backend Switching
	fmt.Println("Testing backend switching...")
	config.Backend = BackendSpeedCrunch
	res1 := wrapper.Evaluate("sin(pi/4)")
	config.Backend = BackendGo
	res2 := wrapper.Evaluate("sin(pi/4)")
	if res1 == res2 {
		fmt.Printf("Switch test FAILED: both backends produced identical string %s\n", res1)
		failed++
	} else {
		fmt.Printf("Switch test PASSED: SC=%s, Go=%s\n", res1, res2)
	}

	if failed > 0 {
		fmt.Printf("\nTotal FAILED tests: %d\n", failed)
		os.Exit(1)
	}
	fmt.Println("\nAll tests PASSED!")
}
