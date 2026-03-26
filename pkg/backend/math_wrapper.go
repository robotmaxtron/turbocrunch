package backend

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v3"
	"github.com/robotmaxtron/turbocrunch/pkg/bridge"
)

// MathBackend defines the type of backend used for calculations.
type MathBackend int

const (
	// BackendSpeedCrunch is the backend that uses the SpeedCrunch core for arbitrary precision calculations.
	BackendSpeedCrunch MathBackend = iota
	// BackendGo is the backend that uses the Go math and cmplx packages for high-performance calculations.
	BackendGo
)

// Theme defines the color scheme for the TUI.
type Theme struct {
	Name            string
	Function        string
	Number          string
	Result          string
	TitleBG         string
	TitleFG         string
	History         string
	HelpKey         string
	HelpDesc        string
	HelpSeparator   string
	ActiveBorder    string
	InactiveBorder  string
	ActiveTitleBG   string
	ActiveTitleFG   string
	InactiveTitleBG string
	InactiveTitleFG string
	Background      string
	Cursor          string
	Bracket         string
	Operator        string
	Nested          string
}

// Config represents the configuration for the evaluator.
type Config struct {
	Backend MathBackend
	Theme   Theme
}

var (
	// Dracula theme
	Dracula = Theme{
		Name:            "Dracula",
		Function:        "#ff79c6", // Pink
		Number:          "#bd93f9", // Purple
		Result:          "#50fa7b", // Green
		TitleBG:         "#6272a4", // Comment
		TitleFG:         "#f8f8f2", // Foreground
		History:         "#6272a4", // Comment
		HelpKey:         "#ffb86c", // Orange
		HelpDesc:        "#6272a4", // Comment
		HelpSeparator:   "#44475a", // Current Line
		ActiveBorder:    "#bd93f9", // Purple
		InactiveBorder:  "#44475a", // Current Line
		ActiveTitleBG:   "#bd93f9", // Purple
		ActiveTitleFG:   "#282a36", // Background
		InactiveTitleBG: "#44475a", // Current Line
		InactiveTitleFG: "#6272a4", // Comment
		Background:      "#282a36",
		Cursor:          "#ff79c6",
		Bracket:         "#f8f8f2", // Foreground
		Operator:        "#8be9fd", // Cyan
		Nested:          "#ffb86c", // Orange
	}

	// Nord theme
	Nord = Theme{
		Name:            "Nord",
		Function:        "#81a1c1", // Frost 3
		Number:          "#b48ead", // Aurora 5
		Result:          "#a3be8c", // Aurora 4
		TitleBG:         "#4c566a", // Polar Night 4
		TitleFG:         "#eceff4", // Snow Storm 3
		History:         "#4c566a", // Polar Night 4
		HelpKey:         "#ebcb8b", // Aurora 3
		HelpDesc:        "#4c566a", // Polar Night 4
		HelpSeparator:   "#3b4252", // Polar Night 2
		ActiveBorder:    "#88c0d0", // Frost 2
		InactiveBorder:  "#3b4252", // Polar Night 2
		ActiveTitleBG:   "#88c0d0", // Frost 2
		ActiveTitleFG:   "#2e3440", // Polar Night 1
		InactiveTitleBG: "#3b4252", // Polar Night 2
		InactiveTitleFG: "#4c566a", // Polar Night 4
		Background:      "#2e3440",
		Cursor:          "#88c0d0",
		Bracket:         "#d8dee9", // Snow Storm 1
		Operator:        "#81a1c1", // Frost 3
		Nested:          "#ebcb8b", // Aurora 3
	}

	// SolarizedDark theme
	SolarizedDark = Theme{
		Name:            "Solarized Dark",
		Function:        "#268bd2", // Blue
		Number:          "#d33682", // Magenta
		Result:          "#859900", // Green
		TitleBG:         "#073642", // Base 02
		TitleFG:         "#93a1a1", // Base 1
		History:         "#586e75", // Base 01
		HelpKey:         "#b58900", // Yellow
		HelpDesc:        "#586e75", // Base 01
		HelpSeparator:   "#073642", // Base 02
		ActiveBorder:    "#268bd2", // Blue
		InactiveBorder:  "#073642", // Base 02
		ActiveTitleBG:   "#268bd2", // Blue
		ActiveTitleFG:   "#002b36", // Base 03
		InactiveTitleBG: "#073642", // Base 02
		InactiveTitleFG: "#586e75", // Base 01
		Background:      "#002b36",
		Cursor:          "#268bd2",
		Bracket:         "#839496", // Base 0
		Operator:        "#2aa198", // Cyan
		Nested:          "#cb4b16", // Orange
	}

	// DefaultTheme (based on original styles)
	DefaultTheme = Theme{
		Name:            "Default",
		Function:        "39",
		Number:          "214",
		Result:          "42",
		TitleBG:         "62",
		TitleFG:         "230",
		History:         "240",
		HelpKey:         "220",
		HelpDesc:        "241",
		HelpSeparator:   "238",
		ActiveBorder:    "62",
		InactiveBorder:  "240",
		ActiveTitleBG:   "62",
		ActiveTitleFG:   "230",
		InactiveTitleBG: "238",
		InactiveTitleFG: "245",
		Background:      "",
		Cursor:          "214",
		Bracket:         "252",
		Operator:        "39",
		Nested:          "208",
	}

	// Miami theme
	Miami = Theme{
		Name:            "Miami",
		Function:        "#ff2f91", // Pink
		Number:          "#00d9f5", // Teal
		Result:          "#50fa7b", // Green
		TitleBG:         "#bd93f9", // Purple
		TitleFG:         "#f8f8f2", // White
		History:         "#6272a4", // Grey
		HelpKey:         "#ffb86c", // Orange
		HelpDesc:        "#6272a4", // Grey
		HelpSeparator:   "#44475a", // Dark Grey
		ActiveBorder:    "#ff2f91", // Pink
		InactiveBorder:  "#44475a", // Dark Grey
		ActiveTitleBG:   "#ff2f91", // Pink
		ActiveTitleFG:   "#242424", // Dark
		InactiveTitleBG: "#44475a", // Dark Grey
		InactiveTitleFG: "#6272a4", // Grey
		Background:      "#1e1e1e", // Dark Background
		Cursor:          "#00d9f5", // Teal
		Bracket:         "#f8f8f2", // White
		Operator:        "#ff2f91", // Pink
		Nested:          "#bd93f9", // Purple
	}
)

// Themes is a list of available themes.
var Themes = []Theme{DefaultTheme, Dracula, Nord, SolarizedDark, Miami}

// EvaluatorWrapper is a high-level wrapper that coordinates between different backends.
type EvaluatorWrapper struct {
	scEvaluator *bridge.Evaluator
	config      *Config
	lastAns     string
}

var scientificConstants = map[string]string{
	"pi":    "3.1415926535897932384626433832795028841971693993751",
	"e":     "2.7182818284590452353602874713526624977572470936999",
	"phi":   "1.6180339887498948482045868343656381177203091798057",
	"light": "299792458",
	"c":     "299792458",
	"G":     "6.67408e-11",
	"h":     "6.626070040e-34",
	"h_bar": "1.054571800e-34",
	"k":     "1.38064852e-23",
	"sigma": "5.670367e-8",
	"N_A":   "6.022140857e23",
	"R":     "8.3144598",
	"F":     "96485.33289",
	"mu_0":  "12.566370614e-7",
	"eps_0": "8.854187817e-12",
	"u":     "1.660539040e-27",
	"g":     "9.80665",
	"au":    "149597870691",
	"ly":    "9.4607304725808e15",
	"pc":    "3.08567802e16",
	"alpha": "7.2973525664e-3",
	"m_e":   "9.10938356e-31",
	"m_p":   "1.672621898e-27",
	"m_n":   "1.674927471e-27",
}

// NewEvaluatorWrapper creates a new instance of EvaluatorWrapper with the provided configuration.
func NewEvaluatorWrapper(config *Config) *EvaluatorWrapper {
	return &EvaluatorWrapper{
		scEvaluator: bridge.NewEvaluator(),
		config:      config,
		lastAns:     "0",
	}
}

// Evaluate takes a mathematical expression and returns the result as a string using the selected backend.
func (e *EvaluatorWrapper) Evaluate(expr string) string {
	expr = strings.TrimSpace(expr)
	if e.config.Backend == BackendSpeedCrunch {
		res, err := e.scEvaluator.EvaluateUpdateAns(expr)
		if err != nil {
			return "Error: " + err.Error()
		}

		if e.scEvaluator.IsUserFunctionAssign() {
			return "Function defined"
		}

		e.lastAns = res
		return res
	}

	res := e.evaluateGo(expr)
	if !strings.HasPrefix(res, "Error:") {
		e.lastAns = res
	}
	return res
}

// SetAngleMode sets the angle mode for the SpeedCrunch backend.
func (e *EvaluatorWrapper) SetAngleMode(mode byte) {
	bridge.SetAngleMode(mode)
}

// GetAngleMode returns the angle mode for the SpeedCrunch backend.
func (e *EvaluatorWrapper) GetAngleMode() byte {
	return bridge.GetAngleMode()
}

// SetResultFormat sets the result format for the SpeedCrunch backend.
func (e *EvaluatorWrapper) SetResultFormat(format byte) {
	bridge.SetResultFormat(format)
}

// GetResultFormat returns the result format for the SpeedCrunch backend.
func (e *EvaluatorWrapper) GetResultFormat() byte {
	return bridge.GetResultFormat()
}

// GetVariable returns the value of a variable from the SpeedCrunch backend.
func (e *EvaluatorWrapper) GetVariable(name string) (string, bool) {
	return e.scEvaluator.GetVariable(name)
}

// Constant represents a mathematical constant from SpeedCrunch.
type Constant = bridge.Constant

// GetConstants returns the list of constants from the SpeedCrunch backend.
func (e *EvaluatorWrapper) GetConstants() []Constant {
	return bridge.GetConstants()
}

// GetUnits returns the list of unit names from the SpeedCrunch backend.
func (e *EvaluatorWrapper) GetUnits() []string {
	return bridge.GetUnits()
}

// FunctionInfo represents metadata for a SpeedCrunch function.
type FunctionInfo = bridge.FunctionInfo

// GetFunctions returns the list of available functions from the SpeedCrunch backend.
func (e *EvaluatorWrapper) GetFunctions() []FunctionInfo {
	return bridge.GetFunctions()
}

// UserFunction represents a user-defined function.
type UserFunction = bridge.UserFunction

// Variable represents a user-defined variable.
type Variable = bridge.Variable

// Function represents a SpeedCrunch function.
type Function = bridge.FunctionInfo

// GetUserFunctions returns the list of user-defined functions from the SpeedCrunch backend.
func (e *EvaluatorWrapper) GetUserFunctions() []UserFunction {
	return e.scEvaluator.GetUserFunctions()
}

// GetVariables returns the list of user-defined variables from the SpeedCrunch backend.
func (e *EvaluatorWrapper) GetVariables() []Variable {
	return e.scEvaluator.GetVariables()
}

// SaveSession saves the current SpeedCrunch session to a file.
func (e *EvaluatorWrapper) SaveSession(filename string) error {
	return e.scEvaluator.SaveSession(filename)
}

// LoadSession loads a SpeedCrunch session from a file.
func (e *EvaluatorWrapper) LoadSession(filename string) error {
	return e.scEvaluator.LoadSession(filename)
}

// UnsetUserFunction removes a user-defined function.
func (e *EvaluatorWrapper) UnsetUserFunction(name string) {
	e.scEvaluator.UnsetUserFunction(name)
}

// UnsetVariable removes a user-defined variable.
func (e *EvaluatorWrapper) UnsetVariable(name string) {
	e.scEvaluator.UnsetVariable(name)
}

// evaluateGo uses the custom parser for high precision or complex arithmetic and nested expressions.
func (e *EvaluatorWrapper) evaluateGo(expr string) string {
	res, err := e.parseAndEvaluate(expr)
	if err != nil {
		return "Error: " + err.Error()
	}

	return res
}

func (e *EvaluatorWrapper) evaluateComplexGo(expr string) string {
	return e.evaluateGo(expr)
}

func parseComplex(s string) (complex128, error) {
	s = strings.ReplaceAll(s, "j", "i")

	// Support simple divisions like pi/2 or 1/2 in complex parser
	if strings.Contains(s, "/") {
		parts := strings.Split(s, "/")
		if len(parts) == 2 {
			p1 := strings.Trim(parts[0], "()")
			p2 := strings.Trim(parts[1], "()")
			v1, err1 := parseComplex(p1)
			v2, err2 := parseComplex(p2)
			if err1 == nil && err2 == nil {
				if v2 == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				return v1 / v2, nil
			}
		}
	}

	// Support simple multiplications like pi*1i
	if strings.Contains(s, "*") {
		parts := strings.Split(s, "*")
		if len(parts) == 2 {
			p1 := strings.Trim(parts[0], "()")
			p2 := strings.Trim(parts[1], "()")
			v1, err1 := parseComplex(p1)
			v2, err2 := parseComplex(p2)
			if err1 == nil && err2 == nil {
				return v1 * v2, nil
			}
		}
	}

	// If it's just 'i' or 'j'
	if s == "i" {
		return 0 + 1i, nil
	}
	// If it ends with i but has no other signs, it might be 2i
	if strings.HasSuffix(s, "i") && !strings.ContainsAny(s[:len(s)-1], "+-") {
		v, err := strconv.ParseFloat(s[:len(s)-1], 64)
		if err == nil {
			return complex(0, v), nil
		}
	}

	// Standard Go complex parser
	return strconv.ParseComplex(s, 128)
}

func (e *EvaluatorWrapper) formatComplex(c complex128) string {
	r, i := real(c), imag(c)
	if math.Abs(i) < 1e-15 {
		// Try to see if we can use high precision for real part if it was a simple constant
		// but c is already float64.
		return fmt.Sprintf("%.25g", r)
	}
	if math.Abs(r) < 1e-15 {
		return fmt.Sprintf("%.25gj", i)
	}
	if i < 0 {
		return fmt.Sprintf("%.25g - %.25gj", r, -i)
	}
	return fmt.Sprintf("%.25g + %.25gj", r, i)
}

func (e *EvaluatorWrapper) formatDecimal(d *apd.Decimal) string {
	return d.Text('f')
}
