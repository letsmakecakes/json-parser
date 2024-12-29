# Go JSON Parser

A high-performance JSON parser implementation in Go, built from scratch.
This project follows a step-by-step approach
to building a fully-featured JSON parser with comprehensive test coverage and validation capabilities.

## Features

- 🚀 High-performance lexical and syntactic analysis
- 🔍 Detailed error reporting with line and column information
- ✅ Full JSON specification support
- 📊 Optional performance benchmarking
- 🔒 Strict mode validation
- 💻 Command-line interface
- 📝 Comprehensive test suite

## Requirements

- Go 1.19 or higher
- Make (optional, for using the Makefile)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/letsmakecakes/json-parser.git
cd json-parser
```

2. Build the project:
```bash
make build
```

Or using Go directly:
```bash
go build -o build/jsonparser ./cmd/parser
```

## Usage

### Basic Usage

Parse a JSON file:
```bash
./build/jsonparser input.json
```

Parse from standard input:
```bash
echo '{"key": "value"}' | ./build/jsonparser -
```

### Command Line Options

```bash
Usage: jsonparser [options] [file]

Options:
  -file string
        JSON file to parse (use '-' for stdin)
  -verbose
        Enable verbose output
  -benchmark
        Show parsing time
  -strict
        Enable strict mode validation
```

### Examples

1. Parse with verbose output:
```bash
./build/jsonparser -verbose input.json
```

2. Parse with benchmarking:
```bash
./build/jsonparser -benchmark input.json
```

3. Parse in strict mode:
```bash
./build/jsonparser -strict input.json
```

## Project Structure

```
.
├── cmd
│   └── parser
│       └── main.go       # Main entry point
├── internal
│   ├── lexer            # Lexical analysis
│   │   ├── lexer.go
│   │   ├── token.go
│   │   └── lexer_test.go
│   ├── parser           # Syntactic analysis
│   │   ├── parser.go
│   │   ├── ast.go
│   │   └── parser_test.go
│   └── validator        # JSON validation
│       ├── validator.go
│       └── validator_test.go
├── pkg
│   └── errors           # Error definitions
│       └── errors.go
├── test                 # Test files
│   ├── step1
│   ├── step2
│   ├── step3
│   ├── step4
│   └── step5
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Development

### Running Tests

Run all tests:
```bash
make test
```

Run specific test steps:
```bash
make test-step1
make test-step2
make test-step3
make test-step4
```

### Development Commands

Format code:
```bash
make fmt
```

Run linter:
```bash
make lint
```

Clean build artifacts:
```bash
make clean
```

## Error Handling

The parser provides detailed error messages with line and column information:

```
Error: 
{"key": "value",}
                ^
unexpected comma at line 1, column 16
```

## Implementation Steps

1. **Step 1**: Basic JSON object parsing (`{}`)
2. **Step 2**: String key-value pairs
3. **Step 3**: Multiple data types (strings, numbers, booleans, null)
4. **Step 4**: Nested objects and arrays
5. **Step 5**: Comprehensive validation and error handling

## Performance

The parser includes optional benchmarking capabilities. When run with the `-benchmark` flag, it provides:
- Total parsing time
- Processing speed (MB/s)

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License—see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by the [JSON specification](https://tools.ietf.org/html/std90)
- Built the following best practices from the "Dragon Book" (Compilers: Principles, Techniques, and Tools)

## Author

Adwaith Rajeev ([@letsmakecakes](https://github.com/letsmakecakes))

## Support

If you have any questions or run into issues, please open an issue in the GitHub repository.