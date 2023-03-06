package evaluator

import (
	"fmt"
	"github.com/jackwilsdon/go-calc/ast"
	"github.com/jackwilsdon/go-calc/token"
	"math"
	"strconv"
)

// op performs a named operation against two values.
func op(a, b float64, op string) (float64, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		return a / b, nil
	case "^":
		return math.Pow(a, b), nil
	default:
		return 0, fmt.Errorf("unsupported operation: %s", op)
	}
}

// Evaluate returns the result of evaluating n with the provided constants.
func Evaluate(n ast.Node, constants map[string]float64) (float64, error) {
	// Evaluate the left and right sides of binary expressions and then
	// perform the described operation on them.
	if b, ok := n.(ast.BinaryExpr); ok {
		left, err := Evaluate(b.Left, constants)

		if err != nil {
			return 0, err
		}

		right, err := Evaluate(b.Right, constants)

		if err != nil {
			return 0, err
		}

		return op(left, right, b.Op)
	}

	// We can interpret the value of a literal as a floating point number.
	if l, ok := n.(ast.Lit); ok {
		switch l.Type {
		case token.NumberToken:
			return strconv.ParseFloat(l.Value, 64)
		case token.ConstantToken:
			key := l.Value
			if key[0] == '+' || key[0] == '-' {
				key = key[1:]
			}
			v, ok := constants[key]
			if !ok {
				return 0, fmt.Errorf("unknown constant %q", key)
			}
			if l.Value[0] == '-' {
				return -v, nil
			}
			return v, nil
		default:
			return 0, fmt.Errorf("unknown literal type %s (%d)", l.Type, l.Type)
		}
	}

	return 0, fmt.Errorf("unknown node %T", n)
}
