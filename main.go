package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/MatiasLyyra/JabaScript/lex"
	"github.com/MatiasLyyra/JabaScript/parser"
)

var ctx = parser.NewContext()
var debug bool

// TODO: Fix - <unary>
func main() {
	// buf := bytes.NewBufferString("-a\nexit\n")
	r := bufio.NewReader(os.Stdin)
	fmt.Println("JabaScript ver 0.2.0 Author: Matias Lyyra")
	fmt.Println("Type help for info")
	for {
		fmt.Print("> ")
		code, err := r.ReadString('\n')
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			continue
		}
		switch code {
		case "exit\n":
			fmt.Println("bye")
			os.Exit(0)
		case "debug\n":
			debug = !debug
		case "help\n":
			fmt.Println(
				`Documentation:
  - Functions:
    - Fn definition: adder = |a b| a + b
    - Fn calling:    adder(1 2)
  - Basic arithmetic:
    - 1 + b, 8 % 3, (a + 2) * 3, etc
  - Data types:
    - Integer
    - Function
Commands:
  - exit
  - debug
    - Enables some debug info
`)
		case "\n":
			continue
		default:
			execute(code)
		}
	}
}

func execute(code string) {
	tokens, err := lex.Tokenize(bytes.NewBufferString(code))
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}
	if debug {
		fmt.Printf("Tokens: %s\n", tokens)
	}
	exp, err := parser.Parse(tokens)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}
	if debug {
		fmt.Printf("Expression: %s\n", exp)
	}
	val, err := exp.Eval(ctx)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}
	fmt.Println(val)
}
