package token_test

import (
	"github.com/jackwilsdon/go-calc/token"
	"strconv"
	"strings"
	"testing"
)

func TestScannerScanAll(t *testing.T) {
	cases := []struct {
		s string
		t []token.Token
	}{
		{
			"1",
			[]token.Token{
				{token.NumberToken, "1", 0},
			},
		},
		{
			"1.0",
			[]token.Token{
				{token.NumberToken, "1.0", 0},
			},
		},
		{
			"0 + 2 / .3",
			[]token.Token{
				{token.NumberToken, "0", 0},
				{token.OperatorToken, "+", 2},
				{token.NumberToken, "2", 4},
				{token.OperatorToken, "/", 6},
				{token.NumberToken, ".3", 8},
			},
		},
		{
			"(1 * 2) + (3 - 4)",
			[]token.Token{
				{token.ParenthesisToken, "(", 0},
				{token.NumberToken, "1", 1},
				{token.OperatorToken, "*", 3},
				{token.NumberToken, "2", 5},
				{token.ParenthesisToken, ")", 6},
				{token.OperatorToken, "+", 8},
				{token.ParenthesisToken, "(", 10},
				{token.NumberToken, "3", 11},
				{token.OperatorToken, "-", 13},
				{token.NumberToken, "4", 15},
				{token.ParenthesisToken, ")", 16},
			},
		},
		{
			"3 + 4 * 2 / (1 - 5) ^ 2 ^ 3",
			[]token.Token{
				{token.NumberToken, "3", 0},
				{token.OperatorToken, "+", 2},
				{token.NumberToken, "4", 4},
				{token.OperatorToken, "*", 6},
				{token.NumberToken, "2", 8},
				{token.OperatorToken, "/", 10},
				{token.ParenthesisToken, "(", 12},
				{token.NumberToken, "1", 13},
				{token.OperatorToken, "-", 15},
				{token.NumberToken, "5", 17},
				{token.ParenthesisToken, ")", 18},
				{token.OperatorToken, "^", 20},
				{token.NumberToken, "2", 22},
				{token.OperatorToken, "^", 24},
				{token.NumberToken, "3", 26},
			},
		},
		{
			"2 * π",
			[]token.Token{
				{token.NumberToken, "2", 0},
				{token.OperatorToken, "*", 2},
				{token.ConstantToken, "π", 4},
			},
		},
		{
			"3.141 / Pi + 2",
			[]token.Token{
				{token.NumberToken, "3.141", 0},
				{token.OperatorToken, "/", 6},
				{token.ConstantToken, "Pi", 8},
				{token.OperatorToken, "+", 11},
				{token.NumberToken, "2", 13},
			},
		},
		{
			"-Pi + -Pi + Pi + Pi",
			[]token.Token{
				{token.OperatorToken, "-", 0},
				{token.ConstantToken, "Pi", 1},
				{token.OperatorToken, "+", 4},
				{token.OperatorToken, "-", 6},
				{token.ConstantToken, "Pi", 7},
				{token.OperatorToken, "+", 10},
				{token.ConstantToken, "Pi", 12},
				{token.OperatorToken, "+", 15},
				{token.ConstantToken, "Pi", 17},
			},
		},
	}

	for i, c := range cases {
		t.Run("case "+strconv.Itoa(i), func(t *testing.T) {
			s, err := token.NewScanner(strings.NewReader(c.s)).ScanAll()
			if err != nil {
				t.Fatal(err)
			}

			if len(s) != len(c.t) {
				t.Fatalf("expected %d tokens but got %d", len(c.t), len(s))
			}

			for j, tok := range s {
				expectedToken := c.t[j]

				if tok.Type != expectedToken.Type {
					t.Errorf("token %d: expected type to be %s (%d) but got %s (%d)", j, expectedToken.Type, expectedToken.Type, tok.Type, tok.Type)
					continue
				}

				if tok.Value != expectedToken.Value {
					t.Errorf("token %d: expected value to be %q but got %q", j, expectedToken.Value, tok.Value)
					continue
				}

				if tok.Position != expectedToken.Position {
					t.Errorf("token %d: expected position to be %d but got %d", j, expectedToken.Position, tok.Position)
					continue
				}
			}
		})
	}
}
