package token

type Type int

const (
	UnknownToken Type = iota
	ParenthesisToken
	NumberToken
	OperatorToken
)

func (t Type) String() string {
	switch t {
	case ParenthesisToken:
		return "Parenthesis"
	case NumberToken:
		return "Number"
	case OperatorToken:
		return "Operator"
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
	return "\"" + t.Value + "\""
}
