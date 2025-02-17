package main

import (
	"fmt"
	"os"

	"github.com/astraikis/harp/internal/interpreter"
	"github.com/astraikis/harp/internal/parser"
	"github.com/astraikis/harp/internal/scanner"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: harp <script>")
	} else {
		runFile(os.Args[1])
	}
}

func runFile(path string) {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error: Unable to find source file.")
		os.Exit(1)
	}

	run(string(file))
}

func run(source string) {
	tokens := scanner.Scan(source)

	stmts, parseErrors := parser.Parse(tokens)
	if parseErrors != nil {
		for _, parseError := range parseErrors {
			fmt.Println(parseError.Error())
		}
	}

	interpreter.Interpret(stmts)
}
