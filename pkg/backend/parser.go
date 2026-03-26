package backend

import (
	"fmt"
	"math/cmplx"
	"strings"
	"unicode"

	"github.com/cockroachdb/apd/v3"
)

// tokenType represents the type of a mathematical token.
type tokenType int

const (
	tokenNumber tokenType = iota
	tokenOperator
	tokenFunction
	tokenLParen
	tokenRParen
	tokenComma
	tokenIdentifier
)

type token struct {
	typ   tokenType
	value string
}

// parser handles the conversion of infix expressions to postfix (RPN)
// and then evaluates the RPN expression.
type parser struct {
	tokens []token
	pos    int
}

func tokenize(expr string) ([]token, error) {
	var tokens []token
	runes := []rune(expr)
	n := len(runes)
	for i := 0; i < n; {
		r := runes[i]
		if unicode.IsSpace(r) {
			i++
			continue
		}

		if unicode.IsDigit(r) || r == '.' {
			start := i
			for i < n && (unicode.IsDigit(runes[i]) || runes[i] == '.' || runes[i] == 'e' || runes[i] == 'E') {
				// Handle scientific notation
				if (runes[i] == 'e' || runes[i] == 'E') && i+1 < n && (runes[i+1] == '+' || runes[i+1] == '-') {
					i += 2
				} else {
					i++
				}
			}
			// Check for 'i' or 'j' suffix for complex numbers
			if i < n && (runes[i] == 'i' || runes[i] == 'j') {
				i++
			}
			tokens = append(tokens, token{tokenNumber, string(runes[start:i])})
		} else if unicode.IsLetter(r) || r == '_' {
			start := i
			for i < n && (unicode.IsLetter(runes[i]) || unicode.IsDigit(runes[i]) || runes[i] == '_') {
				i++
			}
			val := string(runes[start:i])
			if val == "i" || val == "j" {
				tokens = append(tokens, token{tokenNumber, val})
			} else {
				tokens = append(tokens, token{tokenIdentifier, val})
			}
		} else {
			switch r {
			case '+', '-', '*', '/', '^':
				tokens = append(tokens, token{tokenOperator, string(r)})
			case '(':
				tokens = append(tokens, token{tokenLParen, "("})
			case ')':
				tokens = append(tokens, token{tokenRParen, ")"})
			case ',':
				tokens = append(tokens, token{tokenComma, ","})
			default:
				return nil, fmt.Errorf("unexpected character: %c", r)
			}
			i++
		}
	}
	return tokens, nil
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	case "^":
		return 3
	case "unary-":
		return 4
	}
	return 0
}

func isRightAssociative(op string) bool {
	return op == "^"
}

func (e *EvaluatorWrapper) parseAndEvaluate(expr string) (string, error) {
	tokens, err := tokenize(expr)
	if err != nil {
		return "", err
	}

	// Shunting-yard algorithm to convert to RPN
	var output []token
	var stack []token

	for i, t := range tokens {
		switch t.typ {
		case tokenNumber:
			output = append(output, t)
		case tokenIdentifier:
			// Check if next is LParen, then it's a function, else it's a constant
			if i+1 < len(tokens) && tokens[i+1].typ == tokenLParen {
				t.typ = tokenFunction
				stack = append(stack, t)
			} else {
				// Constant replacement
				if val, ok := scientificConstants[t.value]; ok {
					output = append(output, token{tokenNumber, val})
				} else if t.value == "ans" {
					output = append(output, token{tokenNumber, e.lastAns})
				} else {
					return "", fmt.Errorf("unknown identifier: %s", t.value)
				}
			}
		case tokenFunction:
			stack = append(stack, t)
		case tokenComma:
			for len(stack) > 0 && stack[len(stack)-1].typ != tokenLParen {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return "", fmt.Errorf("misplaced comma or mismatched parentheses")
			}
		case tokenOperator:
			op := t.value
			// Handle unary minus
			if op == "-" && (i == 0 || tokens[i-1].typ == tokenLParen || tokens[i-1].typ == tokenOperator || tokens[i-1].typ == tokenComma) {
				op = "unary-"
				t.value = op
			}

			for len(stack) > 0 && stack[len(stack)-1].typ == tokenOperator {
				topOp := stack[len(stack)-1].value
				if (isRightAssociative(op) && precedence(op) < precedence(topOp)) ||
					(!isRightAssociative(op) && precedence(op) <= precedence(topOp)) {
					output = append(output, stack[len(stack)-1])
					stack = stack[:len(stack)-1]
				} else {
					break
				}
			}
			stack = append(stack, t)
		case tokenLParen:
			stack = append(stack, t)
		case tokenRParen:
			for len(stack) > 0 && stack[len(stack)-1].typ != tokenLParen {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return "", fmt.Errorf("mismatched parentheses")
			}
			stack = stack[:len(stack)-1] // pop LParen
			if len(stack) > 0 && stack[len(stack)-1].typ == tokenFunction {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1].typ == tokenLParen {
			return "", fmt.Errorf("mismatched parentheses")
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	// First try evaluating using APD for high precision
	// We only use APD if there are no complex numbers or unsupported functions.
	useAPD := true
	for _, t := range output {
		if t.typ == tokenNumber && (strings.Contains(t.value, "i") || strings.Contains(t.value, "j")) {
			useAPD = false
			break
		}
		if t.typ == tokenFunction {
			switch t.value {
			case "abs", "sqrt", "exp", "ln", "log", "log10", "sin", "cos", "floor", "ceil":
				// supported by APD or our Taylor series
			default:
				useAPD = false
				break
			}
		}
		if !useAPD {
			break
		}
	}

	if useAPD {
		ctx := apd.BaseContext.WithPrecision(50)
		ctx.Traps = 0
		var evalStack []*apd.Decimal
		for _, t := range output {
			switch t.typ {
			case tokenNumber:
				val, _, err := apd.NewFromString(t.value)
				if err != nil {
					useAPD = false
					goto complexFallback
				}
				evalStack = append(evalStack, val)
			case tokenOperator:
				if t.value == "unary-" {
					if len(evalStack) < 1 {
						return "", fmt.Errorf("insufficient operands for unary minus")
					}
					v := evalStack[len(evalStack)-1]
					evalStack[len(evalStack)-1] = new(apd.Decimal).Neg(v)
					continue
				}
				if len(evalStack) < 2 {
					return "", fmt.Errorf("insufficient operands for operator %s", t.value)
				}
				v2 := evalStack[len(evalStack)-1]
				v1 := evalStack[len(evalStack)-2]
				evalStack = evalStack[:len(evalStack)-2]
				res := new(apd.Decimal)
				switch t.value {
				case "+":
					ctx.Add(res, v1, v2)
				case "-":
					ctx.Sub(res, v1, v2)
				case "*":
					ctx.Mul(res, v1, v2)
				case "/":
					if v2.IsZero() {
						return "", fmt.Errorf("division by zero")
					}
					ctx.Quo(res, v1, v2)
				case "^":
					// APD doesn't have Pow for general exponents easily
					// Fallback to complex for power if not integer?
					// For now fallback to complex for ^
					useAPD = false
					goto complexFallback
				}
				evalStack = append(evalStack, res)
			case tokenFunction:
				if len(evalStack) < 1 {
					return "", fmt.Errorf("insufficient operands for function %s", t.value)
				}
				arg := evalStack[len(evalStack)-1]
				evalStack = evalStack[:len(evalStack)-1]
				res := new(apd.Decimal)
				switch t.value {
				case "abs":
					ctx.Abs(res, arg)
				case "sqrt":
					ctx.Sqrt(res, arg)
				case "exp":
					ctx.Exp(res, arg)
				case "ln":
					ctx.Ln(res, arg)
 			case "log", "log10":
 				ctx.Log10(res, arg)
				case "sin":
					// Basic high-precision sin
					term := new(apd.Decimal).Set(arg)
					res.Set(arg)
					x2 := new(apd.Decimal)
					ctx.Mul(x2, arg, arg)
					for i := int64(3); i < 200; i += 2 {
						ctx.Mul(term, term, x2)
						div := new(apd.Decimal).SetInt64(i * (i - 1))
						ctx.Quo(term, term, div)
						term.Neg(term)
						ctx.Add(res, res, term)
						if term.IsZero() {
							break
						}
					}
				case "cos":
					// Basic high-precision cos
					term := new(apd.Decimal).SetInt64(1)
					res.SetInt64(1)
					x2 := new(apd.Decimal)
					ctx.Mul(x2, arg, arg)
					for i := int64(2); i < 200; i += 2 {
						ctx.Mul(term, term, x2)
						div := new(apd.Decimal).SetInt64(i * (i - 1))
						ctx.Quo(term, term, div)
						term.Neg(term)
						ctx.Add(res, res, term)
						if term.IsZero() {
							break
						}
					}
				case "floor":
					ctx.Floor(res, arg)
				case "ceil":
					ctx.Ceil(res, arg)
				default:
					useAPD = false
					goto complexFallback
				}
				evalStack = append(evalStack, res)
			}
		}
		if len(evalStack) == 1 {
			return e.formatDecimal(evalStack[0]), nil
		}
	}

complexFallback:
	// Evaluate RPN using complex128
	var evalStack []complex128
	for _, t := range output {
		switch t.typ {
		case tokenNumber:
			val, err := parseComplex(t.value)
			if err != nil {
				return "", err
			}
			evalStack = append(evalStack, val)
		case tokenOperator:
			if t.value == "unary-" {
				if len(evalStack) < 1 {
					return "", fmt.Errorf("insufficient operands for unary minus")
				}
				v := evalStack[len(evalStack)-1]
				evalStack[len(evalStack)-1] = -v
				continue
			}
			if len(evalStack) < 2 {
				return "", fmt.Errorf("insufficient operands for operator %s", t.value)
			}
			v2 := evalStack[len(evalStack)-1]
			v1 := evalStack[len(evalStack)-2]
			evalStack = evalStack[:len(evalStack)-2]
			var res complex128
			switch t.value {
			case "+":
				res = v1 + v2
			case "-":
				res = v1 - v2
			case "*":
				res = v1 * v2
			case "/":
				if v2 == 0 {
					return "", fmt.Errorf("division by zero")
				}
				res = v1 / v2
			case "^":
				res = cmplx.Pow(v1, v2)
			}
			evalStack = append(evalStack, res)
		case tokenFunction:
			if len(evalStack) < 1 {
				return "", fmt.Errorf("insufficient operands for function %s", t.value)
			}
			arg := evalStack[len(evalStack)-1]
			evalStack = evalStack[:len(evalStack)-1]
			var res complex128
			switch t.value {
			case "abs":
				res = complex(cmplx.Abs(arg), 0)
			case "sin":
				res = cmplx.Sin(arg)
			case "cos":
				res = cmplx.Cos(arg)
			case "tan":
				res = cmplx.Tan(arg)
			case "asin", "arcsin":
				res = cmplx.Asin(arg)
			case "acos", "arccos":
				res = cmplx.Acos(arg)
			case "atan", "arctan":
				res = cmplx.Atan(arg)
			case "asinh", "arsinh":
				res = cmplx.Asinh(arg)
			case "acosh", "arcosh":
				res = cmplx.Acosh(arg)
			case "atanh", "artanh":
				res = cmplx.Atanh(arg)
			case "exp":
				res = cmplx.Exp(arg)
			case "ln":
				res = cmplx.Log(arg)
			case "log", "log10":
				res = cmplx.Log10(arg)
			case "sqrt":
				res = cmplx.Sqrt(arg)
			case "conj":
				res = cmplx.Conj(arg)
			case "phase":
				res = complex(cmplx.Phase(arg), 0)
			case "real":
				res = complex(real(arg), 0)
			case "imag":
				res = complex(imag(arg), 0)
			case "sinh":
				res = cmplx.Sinh(arg)
			case "cosh":
				res = cmplx.Cosh(arg)
			case "tanh":
				res = cmplx.Tanh(arg)
			case "cot":
				res = cmplx.Cot(arg)
			default:
				return "", fmt.Errorf("unknown function: %s", t.value)
			}
			evalStack = append(evalStack, res)
		}
	}

	if len(evalStack) != 1 {
		return "", fmt.Errorf("invalid expression")
	}
	return e.formatComplex(evalStack[0]), nil
}
