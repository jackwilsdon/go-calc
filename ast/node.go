package ast

import (
	"fmt"
)

type Node fmt.Stringer

type Lit string

func (l Lit) String() string {
	return string(l)
}

type BinaryExpr struct {
	Left, Right Node
	Op          string
}

func (b BinaryExpr) String() string {
	return "(" + b.Left.String() + " " + b.Op + " " + b.Right.String() + ")"
}

var _ Node = Lit("")
var _ Node = BinaryExpr{}
