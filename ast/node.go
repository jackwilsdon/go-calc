package ast

import (
	"fmt"
	"github.com/jackwilsdon/go-calc/token"
)

type Node fmt.Stringer

type Lit struct {
	Type  token.Type
	Value string
}

func (l Lit) String() string {
	return l.Value
}

type BinaryExpr struct {
	Left, Right Node
	Op          string
}

func (b BinaryExpr) String() string {
	return "(" + b.Left.String() + " " + b.Op + " " + b.Right.String() + ")"
}

var _ Node = Lit{}
var _ Node = BinaryExpr{}
