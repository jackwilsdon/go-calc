package parser_test

import (
	"github.com/jackwilsdon/go-calc/ast"
	"github.com/jackwilsdon/go-calc/parser"
	"github.com/jackwilsdon/go-calc/token"
	"reflect"
	"strconv"
	"testing"
)

func TestParser(t *testing.T) {
	cases := []struct {
		s string
		n ast.Node
	}{
		{
			"1",
			ast.Lit{Type: token.NumberToken, Value: "1"},
		},
		{
			"1.0",
			ast.Lit{Type: token.NumberToken, Value: "1.0"},
		},
		{
			"0 + 2 / .3",
			ast.BinaryExpr{
				Left: ast.Lit{Type: token.NumberToken, Value: "0"},
				Right: ast.BinaryExpr{
					Left:  ast.Lit{Type: token.NumberToken, Value: "2"},
					Right: ast.Lit{Type: token.NumberToken, Value: ".3"},
					Op:    "/",
				},
				Op: "+",
			},
		},
		{
			"(1 * 2) + (3 - 4)",
			ast.BinaryExpr{
				Left: ast.BinaryExpr{
					Left:  ast.Lit{Type: token.NumberToken, Value: "1"},
					Right: ast.Lit{Type: token.NumberToken, Value: "2"},
					Op:    "*",
				},
				Right: ast.BinaryExpr{
					Left:  ast.Lit{Type: token.NumberToken, Value: "3"},
					Right: ast.Lit{Type: token.NumberToken, Value: "4"},
					Op:    "-",
				},
				Op: "+",
			},
		},
		{
			"3 + 4 * 2 / (1 - 5) ^ 2 ^ 3",
			ast.BinaryExpr{
				Left: ast.Lit{Type: token.NumberToken, Value: "3"},
				Right: ast.BinaryExpr{
					Left: ast.BinaryExpr{
						Left:  ast.Lit{Type: token.NumberToken, Value: "4"},
						Right: ast.Lit{Type: token.NumberToken, Value: "2"},
						Op:    "*",
					},
					Right: ast.BinaryExpr{
						Left: ast.BinaryExpr{
							Left:  ast.Lit{Type: token.NumberToken, Value: "1"},
							Right: ast.Lit{Type: token.NumberToken, Value: "5"},
							Op:    "-",
						},
						Right: ast.BinaryExpr{
							Left:  ast.Lit{Type: token.NumberToken, Value: "2"},
							Right: ast.Lit{Type: token.NumberToken, Value: "3"},
							Op:    "^",
						},
						Op: "^",
					},
					Op: "/",
				},
				Op: "+",
			},
		},
		{
			"2 * π",
			ast.BinaryExpr{
				Left:  ast.Lit{Type: token.NumberToken, Value: "2"},
				Right: ast.Lit{Type: token.ConstantToken, Value: "π"},
				Op:    "*",
			},
		},
		{
			"3.141 / Pi + 2",
			ast.BinaryExpr{
				Left: ast.BinaryExpr{
					Left:  ast.Lit{Type: token.NumberToken, Value: "3.141"},
					Right: ast.Lit{Type: token.ConstantToken, Value: "Pi"},
					Op:    "/",
				},
				Right: ast.Lit{Type: token.NumberToken, Value: "2"},
				Op:    "+",
			},
		},
		{
			"-Pi + -Pi + Pi + Pi",
			ast.BinaryExpr{
				Left: ast.BinaryExpr{
					Left: ast.BinaryExpr{
						Left:  ast.Lit{Type: token.ConstantToken, Value: "-Pi"},
						Right: ast.Lit{Type: token.ConstantToken, Value: "-Pi"},
						Op:    "+",
					},
					Right: ast.Lit{Type: token.ConstantToken, Value: "Pi"},
					Op:    "+",
				},
				Right: ast.Lit{Type: token.ConstantToken, Value: "Pi"},
				Op:    "+",
			},
		},
	}
	for i, c := range cases {
		t.Run("case "+strconv.Itoa(i), func(t *testing.T) {
			n, err := parser.ParseString(c.s)
			if err != nil {
				t.Fatal(err)
			}

			returnedType := reflect.TypeOf(n)
			expectedType := reflect.TypeOf(c.n)
			if returnedType != expectedType {
				t.Fatalf("expected %s but got %s", expectedType.String(), returnedType.String())
			}
			if n != c.n {
				t.Fatalf("expected %q, got %q", c.n, n)
			}
		})
	}
}
