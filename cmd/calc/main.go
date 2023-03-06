package main

import (
	"fmt"
	"github.com/jackwilsdon/go-calc/evaluator"
	"github.com/jackwilsdon/go-calc/parser"
	"math"
	"os"
	"strconv"
	"strings"
)

var constants = map[string]float64{
	"Inf": math.Inf(1),
	"Pi":  math.Pi,
	"π":   math.Pi,
}

func main() {
	args := os.Args[1:]
	var quiet bool
	if len(args) > 0 && args[0] == "-q" {
		args = args[1:]
		quiet = true
	}
	if len(args) == 0 {
		_, _ = fmt.Fprintf(os.Stderr, "usage: %s [-q] sum\n  -q output result only\n", os.Args[0])
		os.Exit(1)
	}

	node, err := parser.ParseString(strings.Join(args, " "))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to parse: %s\n", err)
		os.Exit(1)
	}

	result, err := evaluator.Evaluate(node, constants)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to interpret: %s\n", err)
		os.Exit(1)
	}

	formattedResult := strconv.FormatFloat(result, 'f', -1, 64)
	if quiet {
		fmt.Println(formattedResult)
	} else {
		fmt.Printf("%s = %s\n", node, formattedResult)
	}
}
