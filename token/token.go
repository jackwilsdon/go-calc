package token

type Type int

const (
	UnknownToken Type = iota
	ParenthesisToken
	NumberToken
	OperatorToken
	EOFToken
)

func (t Type) String() string {
	switch t {
	case ParenthesisToken:
		return "Parenthesis"
	case NumberToken:
		return "Number"
	case OperatorToken:
		return "Operator"
	case EOFToken:
		return "EOF"
	default:
		return "Unknown"
	}
}

type Token struct {
	Type     Type
	Value    string
	Position int
}

func (t Token) String() string {
	if t.Type == EOFToken {
		return "EOF"
	}

	return "\"" + t.Value + "\""
}
