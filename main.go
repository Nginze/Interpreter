package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: go run main.go <filename>")
		return
	}

	filename := args[1]

	lexer := NewLexer(&filename)

	lexer.Tokenize()
	lexer.PrintAll()

}
