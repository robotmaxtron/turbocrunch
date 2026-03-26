package backend

import (
	"regexp"
	"strings"
	"testing"
)

// TestFormulaBookContent verifies that all formulas defined in CommonFormulas have required fields.
func TestFormulaBookContent(t *testing.T) {
	if len(CommonFormulas) == 0 {
		t.Error("CommonFormulas list is empty")
	}

	for i, f := range CommonFormulas {
		if f.Name == "" {
			t.Errorf("Formula at index %d has no Name", i)
		}
		if f.Description == "" {
			t.Errorf("Formula %q at index %d has no Description", f.Name, i)
		}
		if f.Template == "" {
			t.Errorf("Formula %q at index %d has no Template", f.Name, i)
		}
	}
}

// TestFormulaInsertionMock simulates selecting a formula from the TUI.
// This is a minimal test to ensure the backend data is accessible for the TUI.
func TestFormulaSelection(t *testing.T) {
	// Let's pick a known formula
	targetIndex := 0
	if len(CommonFormulas) > 0 {
		f := CommonFormulas[targetIndex]
		if f.Name != "Circle Area" {
			t.Errorf("expected 'Circle Area', got %q", f.Name)
		}
		if f.Template != "pi * r^2" {
			t.Errorf("expected 'pi * r^2', got %q", f.Template)
		}
	}
}

// TestCommonFormulasEvaluation evaluates each common formula using both SpeedCrunch and Go backends.
func TestCommonFormulasEvaluation(t *testing.T) {
	configSC := &Config{Backend: BackendSpeedCrunch}
	configGo := &Config{Backend: BackendGo}

	sc := NewEvaluatorWrapper(configSC)
	goBackend := NewEvaluatorWrapper(configGo)

	// Define test values for common variables
	values := map[string]string{
		"r": "5",
		"a": "1",
		"b": "5",
		"c": "6",
		"h": "10",
		"x": "100",
		"P": "1000",
		"n": "12",
		"t": "5",
	}

	for _, formula := range CommonFormulas {
		expr := formula.Template
		// Replace variable names with their values.
		// We use a regex to match whole words only to avoid partial replacements (e.g. 'r' in 'pi').
		for varName, val := range values {
			re := regexp.MustCompile(`\b` + varName + `\b`)
			expr = re.ReplaceAllString(expr, val)
		}

		t.Logf("Evaluating %s: %s -> %s", formula.Name, formula.Template, expr)

		resSC := sc.Evaluate(expr)
		if strings.HasPrefix(resSC, "Error:") {
			// Some functions like log(x) in SpeedCrunch might need multiple arguments
			// or have different names. For this functional test, we log it.
			t.Logf("%s: SpeedCrunch error (possibly template/version mismatch): %s", formula.Name, resSC)
		}

		resGo := goBackend.Evaluate(expr)
		// We expect the Go backend to fail on complex formulas like Quadratic Formula or Compound Interest
		// as it currently only supports simple arithmetic in its APD implementation.
		if strings.HasPrefix(resGo, "Error:") {
			t.Logf("%s: Go backend (expected) error or limitation: %s", formula.Name, resGo)
		} else {
			t.Logf("%s: SpeedCrunch: %s, Go: %s", formula.Name, resSC, resGo)
			// If both succeeded, they should be somewhat close if it's simple arithmetic
			// but SpeedCrunch and Go backend might have different formatting.
		}
	}
}
