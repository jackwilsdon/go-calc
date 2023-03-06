package evaluator

import (
	"fmt"
	"github.com/jackwilsdon/go-calc/ast"
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

// Evaluate returns the result of evaluating n.
func Evaluate(n ast.Node) (float64, error) {
	// Evaluate the left and right sides of binary expressions and then
	// perform the described operation on them.
	if b, ok := n.(ast.BinaryExpr); ok {
		left, err := Evaluate(b.Left)

		if err != nil {
			return 0, err
		}

		right, err := Evaluate(b.Right)

		if err != nil {
			return 0, err
		}

		return op(left, right, b.Op)
	}

	// We can interpret the value of a literal as a floating point number.
	if l, ok := n.(ast.Lit); ok {
		return strconv.ParseFloat(l.String(), 64)
	}

	return 0, fmt.Errorf("unknown node %T", n)
}
