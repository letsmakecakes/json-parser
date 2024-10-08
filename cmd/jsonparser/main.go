package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	filepath := flag.String("file", "", "Path to the JSON fike to parse")
	flag.Parse()

	if *filepath == "" {
		fmt.Println("Please provide a JSON file using -file flag.")
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(*filepath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	// TODO: Call lexer

	// TODO: call Parser
}
