package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/letsmakecakes/jsonparser/internal/lexer"
	"github.com/letsmakecakes/jsonparser/internal/parser"
	"github.com/letsmakecakes/jsonparser/internal/validator"
	e "github.com/letsmakecakes/jsonparser/pkg/errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const maxDepth = 32 // Maximum nesting depth for JSON

type Config struct {
	inputFile  string
	verbose    bool
	benchmark  bool
	strictMode bool
}

func main() {
	config := parseFlags()

	if err := run(config); err != nil {
		_, err2 := fmt.Fprintf(os.Stderr, "Error: %v", err)
		if err2 != nil {
			log.Fatalf("error printing message to console: %v", err2)
		}
		os.Exit(1)
	}

	fmt.Println("✓ JSON is valid")
}

func parseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.inputFile, "file", "", "JSON file to parse (user '-' for stdin)")
	flag.BoolVar(&config.verbose, "verbose", false, "Enable verbose output")
	flag.BoolVar(&config.benchmark, "benchmark", false, "Show paring time")
	flag.BoolVar(&config.strictMode, "strict", false, "Enable strict mode validation")

	flag.Usage = func() {
		_, err := fmt.Fprintf(os.Stderr, "Usage: %s [options] [file]\n\n", filepath.Base(os.Args[0]))
		if err != nil {
			log.Fatalf("error printing message to console: %v", err)
		}
		_, err = fmt.Fprintf(os.Stderr, "Options:")
		if err != nil {
			log.Fatalf("error printing message to console: %v", err)
		}
		flag.PrintDefaults()
	}

	flag.Parse()

	// Check for positional args if a file is not provided as a flag.
	if config.inputFile == "" && flag.NArg() > 0 {
		config.inputFile = flag.Arg(0)
	}

	// Show usage and exit if no file is provided
	if config.inputFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	return config
}

func run(config *Config) error {
	start := time.Now()

	input, err := readInput(config.inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	if config.verbose {
		_, err2 := fmt.Fprintf(os.Stderr, "Parsing %s...\n", getInputName(config.inputFile))
		if err2 != nil {
			log.Fatalf("error printing message to console: %v", err)
		}
	}

	l := lexer.New(string(input))
	p := parser.New(l)
	v := validator.New(maxDepth)

	root, err := p.Parse()
	if err != nil {
		return handleError(input, err)
	}

	if config.strictMode {
		if err := v.Validate(root); err != nil {
			return fmt.Errorf("validation error: %w", err)
		}
	}

	if config.benchmark {
		displayBenchmark(start, len(input))
	}

	return nil
}

func readInput(filename string) ([]byte, error) {
	if filename == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(filename)
}

func getInputName(filename string) string {
	if filename == "-" {
		return "standard input"
	}
	return filepath.Base(filename)
}

func handleError(input []byte, err error) error {
	var parseErr *e.ParseError
	if errors.As(err, &parseErr) {
		return formatParseError(input, parseErr)
	}
	return err
}

func formatParseError(input []byte, err *e.ParseError) error {
	lines := strings.Split(string(input), "\n")
	if err.Line-1 >= len(lines) {
		return fmt.Errorf("error at end of file: %v", err)
	}

	line := lines[err.Line-1]
	pointer := strings.Repeat(" ", err.Column-1) + "^"

	return fmt.Errorf("\n%s\n%s\n%s at line %d, column %d", line, pointer, err.Message, err.Line, err.Column)
}

func displayBenchmark(start time.Time, inputSize int) {
	duration := time.Since(start)
	_, err := fmt.Fprintf(os.Stderr, "Parsing completed in %v\n", duration)
	if err != nil {
		log.Fatalf("error printing message console: %v", err)
	}
	_, err = fmt.Fprintf(os.Stderr, "Processed %.2f MB/s\n", float64(inputSize)/(1024*1024*duration.Seconds()))
	if err != nil {
		log.Fatalf("error printing message to console: %v", err)
	}
}
