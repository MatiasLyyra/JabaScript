package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/MatiasLyyra/JabaScript/lex"
	"github.com/MatiasLyyra/JabaScript/parser"
)

var ctx *parser.Context
var debug bool

// TODO: Fix - <unary>
func main() {
	r := bufio.NewReader(os.Stdin)
	fmt.Println("JabaScript ver 0.4.1 Author: Matias Lyyra")
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
		case "vars\n":
			printVars()
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
			execute(code, false)
		}
	}
}

func init() {
	ctx = parser.NewContext()
	execute("rand = |seed x| || x = (1664525 * x + 1013904223) % 4294967296\n", true)
	execute("fibbonacci = |n| n - 1 ? n ? fibbonacci(n - 1) + fibbonacci(n - 2) : 0 : 1\n", true)
}

func printVars() {
	for id, val := range ctx.Vars {
		fmt.Printf("%s: %s\n", id, val.Val)
	}
}

func execute(code string, silent bool) {
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
	if !silent || debug {
		fmt.Println(val)

	}
}
