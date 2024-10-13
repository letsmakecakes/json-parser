package lexer

import (
	"fmt"
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

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isHighSurrogate(r rune) bool {
	return r >= 0xD800 && r <= 0xDBFF
}

func isLowSurrogate(r rune) bool {
	return r >= 0xDC00 && r <= 0xDFFF
}

func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token

	for {
		l.skipWhitespace() // Skip any whitespace characters

		var tok Token

		switch l.ch {
		case '{':
			tok = Token{Type: TokenLeftBrace, Literal: "{"}
		case '}':
			tok = Token{Type: TokenRightBrace, Literal: "}"}
		case '[':
			tok = Token{Type: TokenLeftBracket, Literal: "["}
		case ']':
			tok = Token{Type: TokenRightBracket, Literal: "]"}
		case ':':
			tok = Token{Type: TokenColon, Literal: ":"}
		case ',':
			tok = Token{Type: TokenComma, Literal: ","} // Create token for comma
		case '"':
			str, err := l.readString()
			if err != nil {
				return nil, err
			}
			tok = Token{Type: TokenString, Literal: str}
			tokens = append(tokens, tok)
			continue
		case 't':
			if l.peekWord(4) == "true" {
				tok = Token{Type: TokenTrue, Literal: "true"}
				l.advanceBy(4)
			} else {
				return nil, NewUnexpectedCharacterError(l.ch)
			}
		case 'f':
			if l.peekWord(5) == "false" {
				tok = Token{Type: TokenFalse, Literal: "false"}
				l.advanceBy(5)
			} else {
				return nil, NewUnexpectedCharacterError(l.ch)
			}
		case 'n':
			if l.peekWord(4) == "null" {
				tok = Token{Type: TokenNull, Literal: "null"}
				l.advanceBy(4)
			} else {
				return nil, NewUnexpectedCharacterError(l.ch)
			}
		case 0: // End of input
			tok = Token{Type: TokenEOF, Literal: ""}
			tokens = append(tokens, tok)
			return tokens, nil
		default:
			if isDigit(l.ch) || l.ch == '-' {
				num, err := l.readNumber()
				if err != nil {
					return nil, err
				}
				tok = Token{Type: TokenNumber, Literal: num}
			} else {
				return nil, NewUnexpectedCharacterError(l.ch)
			}
		}

		tokens = append(tokens, tok) // Append the created token to the tokens slice
		l.readChar()                 // Move to the next character for the next iteration
	}
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readString() (string, error) {
	var strBuilder strings.Builder

	// Read the opening quote
	l.readChar()

	for {
		switch l.ch {
		case '"': // Closing quote found, return the string
			l.readChar() // Move past the closing quote
			return strBuilder.String(), nil
		case '\\': // Handle escape sequences
			l.readChar() // Move past the backslash
			switch l.ch {
			case '"':
				strBuilder.WriteByte('"')
			case '\\':
				strBuilder.WriteByte('\\')
			case '/':
				strBuilder.WriteByte('/')
			case 'b':
				strBuilder.WriteByte('\b')
			case 'f':
				strBuilder.WriteByte('\f')
			case 'n':
				strBuilder.WriteByte('\n')
			case 'r':
				strBuilder.WriteByte('\r')
			case 't':
				strBuilder.WriteByte('\t')
			case 'u':
				// Handle Unicode escape sequences (e.g., \uXXXX)
				runeValue, err := l.readUnicode()
				if err != nil {
					return "", err
				}

				// Check if the rune is a high surrogate
				if isHighSurrogate(runeValue) {
					// Expecting a low surrogate next
					if l.ch != '\\' {
						return "", fmt.Errorf("expected '\\' after high surrogate, got '%c'", l.ch)
					}
					l.readChar() // Skip the backslash
					if l.ch != 'u' {
						return "", fmt.Errorf("expected 'u' after '\\', got '%c'", l.ch)
					}
					l.readChar() // Skip the 'u'

					lowSurrogate, err := l.readUnicode()
					if err != nil {
						return "", err
					}

					if !isLowSurrogate(lowSurrogate) {
						return "", fmt.Errorf("invalid low surrogate: \\u%04X", lowSurrogate)
					}

					// Combine the surrogate pair into a single rune
					combinedRune := utf16.DecodeRune(runeValue, lowSurrogate)
					if combinedRune == utf8.RuneError {
						return "", fmt.Errorf("invalid surrogate pair: \\u%04X\\u%04X", runeValue, lowSurrogate)
					}

					strBuilder.WriteRune(combinedRune)
				} else {
					// Regular Unicode character
					strBuilder.WriteRune(runeValue)
				}
			default:
				return "", fmt.Errorf("unexpected character: '%c' in string escape", l.ch)
			}
		case 0: // End of input, but no closing quote found
			return "", fmt.Errorf("unterminated string")
		default:
			strBuilder.WriteByte(l.ch)
		}

		l.readChar() // Read the next character
	}
}

func (l *Lexer) readUnicode() (rune, error) {
	var hex string
	for i := 0; i < 4; i++ {
		if !isHexDigit(l.ch) {
			return 0, fmt.Errorf("invalid Unicode escape sequence: \\u%s", hex)
		}
		hex += string(l.ch)
		l.readChar()
	}

	var unicodeValue uint32
	_, err := fmt.Sscanf(hex, "%04x", &unicodeValue)
	if err != nil {
		return 0, fmt.Errorf("invalid Unicode escape sequence: \\u%s", hex)
	}

	return rune(unicodeValue), nil
}

func isHexDigit(ch byte) bool {
	return ('0' <= ch && ch <= '9') || ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
}

func (l *Lexer) readNumber() (string, error) {
	position := l.position

	if l.ch == '-' {
		l.readChar()
	}

	hasDigits := false
	for isDigit(l.ch) {
		hasDigits = true
		l.readChar()
	}

	if l.ch == '.' {
		l.readChar()

		for isDigit(l.ch) {
			hasDigits = true
			l.readChar()
		}
	}

	if !hasDigits {
		return "", fmt.Errorf("invalid number format")
	}

	return l.input[position:l.position], nil
}

func (l *Lexer) peekWord(n int) string {
	end := l.readPosition + n - 1
	if end > len(l.input) {
		end = len(l.input)
	}
	return l.input[l.position:end]
}

func (l *Lexer) advanceBy(n int) {
	for i := 0; i < n; i++ {
		l.readChar()
	}
}
