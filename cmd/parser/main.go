package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/letsmakecakes/jsonparser/internal/lexer"
	"github.com/letsmakecakes/jsonparser/internal/parser"
)

func main() {
	filepath := flag.String("file", "", "Path to the JSON fike to parse")
	flag.Parse()

	if *filepath == "" {
		fmt.Println("Please provide a JSON file using -file flag.")
		os.Exit(1)
	}

	data, err := os.ReadFile(*filepath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	lex := lexer.NewLexer(string(data))
	tokens, lexErr := lex.Tokenize()
	if lexErr != nil {
		fmt.Println("Lexing Error:", lexErr)
		os.Exit(1)
	}

	_, parseErr := parser.Parse(tokens)
	if parseErr != nil {
		fmt.Println("Parsing Error:", parseErr)
		os.Exit(1)
	}

	fmt.Println("Valid JSON")
	os.Exit(0)
}
