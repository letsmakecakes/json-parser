package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/letsmakecakes/jsonparser/internal/lexer"
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

	// TODO: call Parser
}
