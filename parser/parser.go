package parser

import (
	"errors"
	"fmt"
	"github.com/jackwilsdon/go-calc/ast"
	"github.com/jackwilsdon/go-calc/token"
	"io"
	"strings"
)

const (
	leftAssociativity = iota
	rightAssociativity
)

var operators = map[string]struct {
	precedence, associativity int
}{
	"+": {1, leftAssociativity},
	"-": {1, leftAssociativity},
	"*": {2, leftAssociativity},
	"/": {2, leftAssociativity},
	"^": {3, rightAssociativity},
}

// collapseSigns converts a sequence of signs into a single sign.
// For example, "+-+-+-" is converted into just "-".
func collapseSigns(signs string) string {
	positive := true

	for _, s := range signs {
		// Flip the sign if we see a negative.
		if s == '-' {
			positive = !positive
		}
	}

	if positive {
		return "+"
	}

	return "-"
}

type parser struct {
	scanner   *token.Scanner
	nextToken *token.Token
}

// peek returns the next token without advancing past it.
func (p *parser) peek() (token.Token, error) {
	// If we already know what the next token is then we can just return it.
	if p.nextToken != nil {
		return *p.nextToken, nil
	}

	// Otherwise fetch the next token from the scanner.
	nextToken, err := p.scanner.Scan()

	if err != nil {
		return token.Token{}, err
	}

	// Store the next token for future peek calls.
	p.nextToken = &nextToken

	return nextToken, nil
}

// next returns the next token and advances past it.
func (p *parser) next() (token.Token, error) {
	// If there's no next token then we need to scan for one.
	if p.nextToken == nil {
		return p.scanner.Scan()
	}

	// We already know what the next token is so clear the next token and
	// return it.
	nextToken := p.nextToken
	p.nextToken = nil
	return *nextToken, nil
}

func (p *parser) expression(minimumPrecedence int) (ast.Node, error) {
	left, err := p.factor()

	if err != nil {
		return nil, err
	}

	for {
		t, err := p.peek()

		if err == io.EOF {
			// End of expression.
			break
		} else if err != nil {
			return nil, err
		}

		// Stop now if it's not an operator as we have reached the end of the
		// expression.
		if t.Type != token.OperatorToken {
			break
		}

		// Look up some information about the operator.
		op, valid := operators[t.Value]

		if !valid {
			return nil, fmt.Errorf("unknown operator %s at %d", t, t.Position)
		}

		// If this operator has a lower precedence than the minimum then we
		// don't want to look at it now.
		if op.precedence < minimumPrecedence {
			break
		}

		nextMinimumPrecedence := op.precedence

		// If the operator has left associativity then the next minimum
		// precedence is the precedence of the current operator plus one.
		if op.associativity == leftAssociativity {
			nextMinimumPrecedence += 1
		}

		// Consume the token now that we're sure everything is good.
		if _, err := p.next(); err != nil {
			return nil, err
		}

		// Recursively work out the right hand side of the expression.
		right, err := p.expression(nextMinimumPrecedence)

		if err != nil {
			return nil, err
		}

		// Set the left hand side to the newly generated binary expression and
		// go around again.
		left = ast.BinaryExpr{Left: left, Right: right, Op: t.Value}
	}

	return left, nil
}

func (p *parser) factor() (ast.Node, error) {
	t, err := p.next()
	if err == io.EOF {
		return nil, errors.New("unexpected EOF, expected a factor")
	}

	if err != nil {
		return nil, err
	}

	// Handle unary prefixes.
	if t.Type == token.OperatorToken && (t.Value == "+" || t.Value == "-") {
		signs := t.Value

		// Keep consuming operators and adding them onto the sign string.
		for t.Type == token.OperatorToken && (t.Value == "+" || t.Value == "-") {
			t, err = p.next()
			if err == io.EOF {
				return nil, errors.New("unexpected EOF, expected a number")
			} else if err != nil {
				return nil, err
			}

			signs += t.Value
		}

		// We require a number to attach the signs to.
		if t.Type != token.NumberToken {
			return nil, fmt.Errorf("unexpected %s, expected a number at %d", t, t.Position)
		}

		// Collapse the signs and attach them to the value to be parsed.
		return ast.Lit(collapseSigns(signs) + t.Value), nil
	}

	// Numbers are just literal values.
	if t.Type == token.NumberToken {
		return ast.Lit(t.Value), nil
	}

	// Handle expressions in parentheses.
	if t.Type == token.ParenthesisToken && t.Value == "(" {
		// Evaluate the expression after the opening parenthesis.
		expr, err := p.expression(1)

		if err != nil {
			return nil, err
		}

		t, err = p.next()

		if err == io.EOF {
			return nil, errors.New("unexpected EOF, expected closing parenthesis")
		} else if err != nil {
			return nil, err
		}

		// We expect a closing parenthesis now, as we've already evaluated the
		// inner expression.
		if t.Type != token.ParenthesisToken || t.Value != ")" {
			return nil, fmt.Errorf("unexpected %s, expected closing parenthesis at %d", t, t.Position)
		}

		return expr, nil
	}

	return nil, fmt.Errorf("unexpected %s, expected a factor at %d", t, t.Position)
}

func ParseScanner(s *token.Scanner) (ast.Node, error) {
	p := parser{scanner: s}

	node, err := p.expression(1)
	if err != nil {
		return nil, err
	}

	t, err := p.peek()
	if err == io.EOF {
		// No trailing tokens.
		return node, nil
	} else if err != nil {
		return nil, err
	}

	// There are more tokens after the expression.
	return nil, fmt.Errorf("unexpected trailing %s at %d", t, t.Position)
}

func ParseReader(r io.Reader) (ast.Node, error) {
	return ParseScanner(token.NewScanner(r))
}

func ParseString(s string) (ast.Node, error) {
	return ParseReader(strings.NewReader(s))
}
