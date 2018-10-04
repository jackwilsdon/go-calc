package token

import (
	"bufio"
	"io"
	"strings"
)

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

func isDigit(r rune) bool {
	return (r >= '0' && r <= '9') || r == '.'
}

func isOperator(r rune) bool {
	return r == '+' || r == '-' || r == '*' || r == '/' || r == '^'
}

func isParenthesis(r rune) bool {
	return r == '(' || r == ')'
}

func handleError(position int, err error) (Token, error) {
	if err == io.EOF {
		return Token{Type: EOFToken, Position: position}, nil
	}

	return Token{}, err
}

// Scanner converts a stream of runes into a stream of tokens.
type Scanner struct {
	r        *bufio.Reader
	position int
}

// read reads and returns the next rune.
func (s *Scanner) read() (rune, error) {
	r, _, err := s.r.ReadRune()

	if err == nil {
		s.position++
	}

	return r, err
}

// unread unreads the last read rune.
func (s *Scanner) unread() error {
	err := s.r.UnreadRune()

	if err == nil {
		s.position--
	}

	return err
}

// scan reads the next rune if check returns true for the rune.
func (s *Scanner) scan(check func(r rune) bool) (rune, bool, error) {
	// Read the next rune.
	r, err := s.read()

	if err != nil {
		return 0, false, err
	}

	// If the rune is valid then return it.
	if check(r) {
		return r, true, nil
	}

	// Unread the rune.
	return 0, false, s.unread()
}

// scanWhile reads runes until check returns false for a rune.
func (s *Scanner) scanWhile(check func(r rune) bool) (string, error) {
	b := strings.Builder{}

	for {
		// Read the next rune.
		r, err := s.read()

		// If it's an EOF then stop.
		if err == io.EOF {
			return b.String(), nil
		}

		if err != nil {
			return "", err
		}

		// If the rune doesn't match the check then unread it and stop.
		if !check(r) {
			if err := s.unread(); err != nil {
				return "", err
			}

			return b.String(), nil
		}

		// Store the read rune as it matches the check.
		b.WriteRune(r)
	}
}

// Scan reads and returns the next token.
func (s *Scanner) Scan() (Token, error) {
	var startPosition int

	// Eat up all the whitespace as we don't really care about it.
	if _, err := s.scanWhile(isWhitespace); err != nil {
		return handleError(s.position, err)
	}

	// Update the start position before we scan for digits.
	startPosition = s.position

	// Eat as many digits as we can.
	digit, err := s.scanWhile(isDigit)

	if err != nil {
		return handleError(s.position, err)
	}

	// If there's any digits then it's a number.
	if len(digit) > 0 {
		return Token{NumberToken, digit, startPosition}, nil
	}

	// Update the start position before we scan for an operator.
	startPosition = s.position

	// Just try and eat one operator.
	op, isOp, err := s.scan(isOperator)

	if err != nil {
		return handleError(s.position, err)
	}

	// If it did match then it's an operator.
	if isOp {
		return Token{OperatorToken, string(op), startPosition}, nil
	}

	// Update the start position before we scan for parentheses.
	startPosition = s.position

	// Just try and eat one parenthesis.
	paren, isParen, err := s.scan(isParenthesis)

	if err != nil {
		return handleError(s.position, err)
	}

	// If it did match then it's a parenthesis.
	if isParen {
		return Token{ParenthesisToken, string(paren), startPosition}, nil
	}

	// We don't know what type it is, just move along one rune anyway.
	if _, err := s.read(); err != nil {
		return handleError(s.position, err)
	}

	// This is an unknown type.
	return Token{Type: UnknownToken}, nil
}

// ScanAll reads and returns all tokens until io.EOF is returned by the
// underlying reader.
func (s *Scanner) ScanAll() ([]Token, error) {
	var ts []Token

	for {
		t, err := s.Scan()

		if err != nil {
			return nil, err
		}

		// If it's an EOF then don't add it to the array and just return it.
		if t.Type == EOFToken {
			return ts, nil
		}

		// Add the token to the array.
		ts = append(ts, t)
	}
}

// NewScanner creates a new scanner which reads from r.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r), position: 0}
}
