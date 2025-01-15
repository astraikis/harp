package main

import (
	"fmt"
	"os"
)

var hadError bool = false

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: harp <script>")
	} else {
		runFile(os.Args[1])
	}

	if hadError {
		fmt.Println("Exited with error.")
	} else {
		fmt.Println("Exited without error.")
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
	// Scan
	// fmt.Printf("Scanning...\n\n")
	tokens := scan(source)
	// printTokens()
	// Parse
	// fmt.Printf("Parsing %d tokens...\n\n", len(tokens))
	statements := parse(tokens)
	// printStatements()
	// Interpret
	interpret(statements)
}

func harpError(lexeme string, message string, column int, line int) {
	report(lexeme, message, column, line)
	hadError = true
}

func report(lexeme string, message string, column int, line int) {
	fmt.Println("Error: " + message)
	fmt.Printf("   [line %d]\n\n", line)
	fmt.Println("   " + lexeme)
	fmt.Printf("   ^\n\n")
}
