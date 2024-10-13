# JSON Parser in Go

This is a lightweight, hand-written JSON parser implemented in Go. It lexes and parses JSON input, supporting standard JSON data types like objects, arrays, strings, numbers, booleans, and `null`. The parser features a lookahead (`peek`) mechanism for efficient and error-resilient parsing.

## Features

- **Lexer and Parser**: Converts JSON input into a Go data structure.
- **Supports JSON Data Types**:
  - Objects (`{...}`)
  - Arrays (`[...]`)
  - Strings
  - Numbers (integers and floating point)
  - Booleans (`true`/`false`)
  - Null (`null`)
- **Error Handling**: Captures and reports detailed parsing errors.
- **Lookahead Support**: Uses `peek` functionality for lookahead during parsing.
- **Unit Tests**: Comprehensive tests to ensure robustness.

## Installation

To get started, clone the repository and run the parser:

```bash
git clone https://github.com/your-username/json-parser-go.git
cd json-parser-go
go mod tidy
```

### Lexer

The lexer scans the JSON input and breaks it into tokens. Each token has a type (e.g., string, number, left brace) and a literal value.

### Parser

The parser converts tokens into corresponding Go data structures. It supports objects, arrays, and primitive types, including lookahead functionality with a `peek` mechanism for efficient parsing.

## Contributing

Contributions are welcome! If you'd like to improve the parser, please fork the repository and create a pull request with your changes.

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
