package lexer

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

// TokenType defines the type of lexical tokens
type TokenType string

// Token types
const (
	TokenLeftBrace    TokenType = "{"
	TokenRightBrace   TokenType = "}"
	TokenLeftBracket  TokenType = "["
	TokenRightBracket TokenType = "]"
	TokenColon        TokenType = ":"
	TokenComma        TokenType = ","
	TokenString       TokenType = "STRING"
	TokenNumber       TokenType = "NUMBER"
	TokenTrue         TokenType = "TRUE"
	TokenFalse        TokenType = "FALSE"
	TokenNull         TokenType = "NULL"
	TokenEOF          TokenType = "EOF"
)

// Token represents a lexical token with type and literal value
type Token struct {
	Type    TokenType
	Literal string
	Line    int // Line number in input
	Column  int // Column number in input
}

// Lexer represents a lexical scanner
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           rune // current char under examination
	line         int  // current line number
	column       int  // current column number
}

// NewLexer initializes a new Lexer with the given input
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// readChar reads the next character and updates positions
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		r, size := utf8.DecodeLastRuneInString(l.input[l.readPosition:])
		l.ch = r
		l.readPosition += size
		l.position = l.readPosition
		l.column++
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
	}
	l.position = l.readPosition
	l.readPosition++
}

// peekChar peeks ahead to the next character without advancing the lexer
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeLastRuneInString(l.input[l.readPosition:])
	return r
}

// skipWhiteSpace skips over any whitespace characters
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

// Tokenize converts the input string into a slice of Tokens
func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token

	for l.ch != 0 {
		l.skipWhitespace() // Skip any whitespace characters

		var tok Token
		tok.Line = l.line
		tok.Column = l.column

		switch l.ch {
		case '{':
			tok = Token{Type: TokenLeftBrace, Literal: "{", Line: l.line, Column: l.column}
		case '}':
			tok = Token{Type: TokenRightBrace, Literal: "}", Line: l.line, Column: l.column}
		case '[':
			tok = Token{Type: TokenLeftBracket, Literal: "[", Line: l.line, Column: l.column}
		case ']':
			tok = Token{Type: TokenRightBracket, Literal: "]", Line: l.line, Column: l.column}
		case ':':
			tok = Token{Type: TokenColon, Literal: ":", Line: l.line, Column: l.column}
		case ',':
			tok = Token{Type: TokenComma, Literal: ",", Line: l.line, Column: l.column} // Create token for comma
		case '"':
			str, err := l.readString()
			if err != nil {
				return nil, err
			}
			tok = Token{Type: TokenString, Literal: str, Line: l.line, Column: l.column}
			tokens = append(tokens, tok)
			continue
		case 't':
			if l.peekKeyWord("true") {
				tok = Token{Type: TokenTrue, Literal: "true", Line: l.line, Column: l.column}
				l.advanceBy(len("true"))
			} else {
				return nil, fmt.Errorf("Lexer error at line %d, column %d: invalid token starting with 't'", l.line, l.column)
			}
		case 'f':
			if l.peekKeyWord("false") {
				tok = Token{Type: TokenFalse, Literal: "false", Line: l.line, Column: l.column}
				l.advanceBy(len("false"))
			} else {
				return nil, fmt.Errorf("Lexer error at line %d, column %d: invalid token starting with 'f'", l.line, l.column)
			}
		case 'n':
			if l.peekKeyWord("null") {
				tok = Token{Type: TokenNull, Literal: "null", Line: l.line, Column: l.column}
				l.advanceBy(len("null"))
			} else {
				return nil, fmt.Errorf("Lexer error at line %d, column %d: invalid token starting with 'n'", l.line, l.column)
			}
		default:
			if l.isStartOfNumber(l.ch) {
				num, err := l.readNumber()
				if err != nil {
					return nil, fmt.Errorf("Lexer error at line %d, column %d: %v", l.line, l.column, err)
				}
				tok = Token{Type: TokenNumber, Literal: num, Line: l.line, Column: l.column}
			} else {
				return nil, fmt.Errorf("Lexer error at line %d, column %d: unexpected character: %c", l.ch, l.line, l.column)
			}
		}

		tokens = append(tokens, tok) // Append the created token to the tokens slice
		l.readChar()                 // Move to the next character for the next iteration
	}

	// Append EOF token
	tokens = append(tokens, Token{Type: TokenEOF, Literal: "", Line: l.line, Column: l.column})

	return tokens, nil
}

// peekKeyword checks if the upcoming characters match the expected keyword
func (l *Lexer) peekKeyWord(expected string) bool {
	end := l.readPosition + len(expected)
	if end > len(l.input) {
		return false
	}

	return l.input[l.readPosition:end] == expected
}

// advanceBy advances the lexer by n characters
func (l *Lexer) advanceBy(n int) {
	for i := 0; i < n; i++ {
		l.readChar()
	}
}

// isStartOfNumber checks if the rune can start a number
func (l *Lexer) isStartOfNumber(r rune) bool {
	return r == '-' || unicode.IsDigit(r)
}

// readNumber reads a number token from the input, including exponents
func (l *Lexer) readNumber() (string, error) {
	startPos := l.position

	if err := l.consumeMinus(); err != nil {
		return "", err
	}

	if err := l.consumeInteger(); err != nil {
		return "", err
	}

	if err := l.consumeFraction(); err != nil {
		return "", err
	}

	if err := l.consumeExponent(); err != nil {
		return "", err
	}

	numStr := l.input[startPos:l.position]

	// Validate number using strconv
	if _, err := strconv.ParseFloat(numStr, 64); err != nil {
		return "", fmt.Errorf("invalid number format: %v", err)
	}

	// Ensurr that the number is not followed by a letter or digit
	if unicode.IsLetter(l.ch) || isDigit(l.ch) {
		return "", fmt.Errorf("invalid character following number")
	}

	return numStr, nil
}

// consumeMinus handles the optional minus sign
func (l *Lexer) consumeMinus() error {
	if l.ch == '-' {
		l.readChar()
	}
	return nil
}

// consumeInteger parses the integer part and enforces no leading zeros
func (l *Lexer) consumeInteger() error {
	if l.ch == '0' {
		l.readChar()
		// Leading zeros are not allowed unless the number is exactly '0'
		if isDigit(l.ch) {
			return fmt.Errorf("invalid number format: leading zeros are not allowed")
		}
	} else if isDigitOneToNine(l.ch) {
		for isDigit(l.ch) {
			l.readChar()
		}
	} else {
		return fmt.Errorf("expected digit in number")
	}

	return nil
}

// consumeFraction parses the fractional part of the number
func (l *Lexer) consumeFraction() error {
	if l.ch == '.' {
		l.readChar()
		if !isDigit(l.ch) {
			return fmt.Errorf("expected digit after decimal point")
		}
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return nil
}

// consumeExponent parses the exponent part of the number
func (l *Lexer) consumeExponent() error {
	if l.ch == 'e' || l.ch == 'E' {
		l.readChar()
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}
		if !isDigit(l.ch) {
			return fmt.Errorf("expected digit after exponent")
		}
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return nil
}

// isDigit checks if the rune is a digit (0-9)
func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// isDigitOnetoNine checks if the rune is a digit between 1 and 9
func isDigitOneToNine(ch rune) bool {
	return ch >= '1' && ch <= '9'
}

// readString reads a string token, handling escape sequences and Unicode
func (l *Lexer) readString() (string, error) {
	var strBuilder strings.Builder

	l.readChar() // Skip the opening quote

	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case '"':
				strBuilder.WriteRune('"')
			case '\\':
				strBuilder.WriteRune('\\')
			case '/':
				strBuilder.WriteRune('/')
			case 'b':
				strBuilder.WriteRune('\b')
			case 'f':
				strBuilder.WriteRune('\f')
			case 'n':
				strBuilder.WriteRune('\n')
			case 'r':
				strBuilder.WriteRune('\r')
			case 't':
				strBuilder.WriteRune('\t')
			case 'u':
				// Handle Unicode escape sequence
				r, err := l.readUnicode()
				if err != nil {
					return "", err
				}
				strBuilder.WriteRune(r)
			default:
				return "", fmt.Errorf("invalid escape character: '\\%c'", l.ch)
			}
		} else {
			strBuilder.WriteRune(l.ch)
		}
		l.readChar()
	}

	if l.ch != '"' {
		return "", fmt.Errorf("unterminated string literal")
	}

	return strBuilder.String(), nil
}

// readUnicode reads a Unicode escape sequence and returns the corresponding rune
func (l *Lexer) readUnicode() (rune, error) {
	var hexDigits [4]rune
	for i := 0; i < 4; i++ {
		l.readChar()
		if !isHexDigit(l.ch) {
			return 0, fmt.Errorf("invalid unicode escape sequence")
		}
		hexDigits[i] = l.ch
	}

	codePoint := hexToInt(hexDigits)
	r := rune(codePoint)

	if unicode.IsSurrogate(r) {
		if !l.peekUnicodeSurrogatePair() {
			return 0, fmt.Errorf("invalid surrogate pair in Unicode escape")
		}
		// Read the low surrogate
		l.readChar() // Skip 'u'
		l.readChar() // Move to low surrogate
		lexHexDigits, err := l.readUnicodeSurrogate()
		if err != nil {
			return 0, err
		}
		r = utf16.DecodeRune(r, lexHexDigits)
		if r == utf8.RuneError {
			return 0, fmt.Errorf("invalid surrogate pair")
		}
	}

	return r, nil
}

// peekUnicodeSurrogatePair checks if the next sequence is a low surrogate pair
func (l *Lexer) peekUnicodeSurrogatePair() bool {
	backupPosition := l.readPosition
	backupLine := l.line
	backupColumn := l.column

	// Expecting '\', 'u' followed by four hex digits
	if l.peekChar() != '\\' {
		return false
	}
	l.readChar()
	if l.peekChar() != 'u' {
		l.readPosition = backupPosition
		l.line = backupLine
		l.column = backupColumn
		return false
	}
	for i := 0; i < 5; i++ { // Skip '\', 'u', and four hex digits
		l.readChar()
		if i < 2 && (i == 0 && l.ch != '\\' || i == 1 && l.ch != 'u') {
			l.readPosition = backupPosition
			l.line = backupLine
			l.column = backupColumn
			return false
		}
		if i >= 2 && !isHexDigit(l.ch) {
			l.readPosition = backupPosition
			l.line = backupLine
			l.column = backupColumn
			return false
		}
	}
	l.readPosition = backupPosition
	l.line = backupLine
	l.column = backupColumn
	return true
}

// readUnicodeSurrogate reads the low surrogate after '\u'
func (l *Lexer) readUnicodeSurrogate() (rune, error) {
	var hexDigits [4]rune
	for i := 0; i < 4; i++ {
		l.readChar()
		if !isHexDigit(l.ch) {
			return 0, fmt.Errorf("invalid Unicode low surrogate escape sequence")
		}
		hexDigits[i] = l.ch
	}

	codePoint := hexToInt(hexDigits)
	r := rune(codePoint)
	if !unicode.IsLowSurrogate(r) {
		return 0, fmt.Errorf("invalid low surrogate in Unicode escape")
	}

	return r, nil
}

// isHexDigit checks if the rune is a valid hexadecimal digit
func isHexDigit(r rune) bool {
	return ('0' <= r && r <= '9') || ('a' <= r && r <= 'f') || ('A' <= r && r <= 'F')
}

// hexToInt converts an array of four hex runes to an integer
func hexToInt(hex [4]rune) int {
	var value int
	for _, r := range hex {
		value <<= 4
		switch {
		case r >= '0' && r <= '9':
			value += int(r - '0')
		case r >= 'a' && r <= 'f':
			value += int(r-'a') + 10
		case r >= 'A' && r <= 'F':
			value += int(r-'A') + 10
		}
	}
	return value
}
